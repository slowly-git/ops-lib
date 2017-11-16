#!/usr/bin/env python
#-*-coding: utf-8 -*-

import sys
import time
import datetime
import optparse
import re
import os
import pymongo
import bson.son as son
from zabbixConfig import *

dbStatus = {}

def check_ismaster(con):
	try:
		data = con.admin.command('isMaster')
	except Excetion, e:
		print e
	else:
		return data['ismaster']

def get_server_status(con):
	try:
		if not check_ismaster(con):
			set_read_preference(con.admin)
		data = con.admin.command(pymongo.son_manipulator.SON([('serverStatus', 1)]))
	except:
		data = con.admin.command(son.SON([('serverStatus', 1)]))
	return data

def mongo_connect(host=None, port=None):
	try:
		con = pymongo.MongoClient(host, port)
	except Exception, e:
		return 1,
	return 0, con

def exit_with_general_critical(e):
    if isinstance(e, SystemExit):
        return e
    else:
        print "CRITICAL - General MongoDB Error:", e
    return 2    

def set_read_preference(db):
	db.read_preference = pymongo.ReadPreference.SECONDARY
	
def check_connections(data):
	try:
		dbStatus['mongo.connections.current'] = float(data['connections']['current'])
		dbStatus['mongo.connections.available'] = float(data['connections']['available'])
		dbStatus['mongo.connections.used_percent'] =  int(float(dbStatus['mongo.connections.current'] / (dbStatus['mongo.connections.available'] + dbStatus['mongo.connections.current'])) * 100)
		dbStatus['mongo.connections.total'] = dbStatus['mongo.connections.current'] + dbStatus['mongo.connections.available']
		
	except Exception, e:
		print exit_with_general_critical(e)

def check_memory(data):
	try:
		if not data['mem']['supported'] and not mapped_memory:
			print "OK - Platform not supported for memory info"
		try:
			dbStatus['mongo.mem.mem_resident'] = float(data['mem']['resident']) / 1024.0
		except:
			dbStatus['mongo.mem.mem_resident'] = 0
		try:
			dbStatus['mongo.mem.mem_virtual'] = float(data['mem']['virtual']) / 1024.0
		except:
			dbStatus['mongo.mem.mem_virtual'] = 0
		try:
			dbStatus['mongo.mem.mem_mapped'] = float(data['mem']['mapped']) / 1024.0
		except:
			dbStatus['mongo.mem.mem_mapped'] = 0
		try:
			dbStatus['mongo.mem.mem_mapped_journal'] = float(data['mem']['mappedWithJournal']) / 1024.0
		except:
			dbStatus['mongo.mem.mem_mapped_journal'] = 0
	except Exception, e:
		print e

def check_lock(data):
	try:
		dbStatus['mongo.lock.lock_percentage'] = float(data['globalLock']['lockTime']) / float(data['globalLock']['totalTime']) * 100
		
	except Exception, e:
		print exit_with_general_critical(e)

def check_flushing(data):
	try:
		dbStatus['mongo.flushing.average_ms'] = float(data['backgroundFlushing']['average_ms'])
		dbStatus['mongo.flushing.last_ms'] = float(data['backgroundFlushing']['last_ms'])
		dbStatus['mongo.flushing.flushes'] = float(data['backgroundFlushing']['flushes'])
		dbStatus['mongo.flushing.total_ms'] = float(data['backgroundFlushing']['total_ms'])
		dbStatus['mongo.flushing.last_finished'] = data['backgroundFlushing']['last_finished'].strftime('%Y-%m-%d %H:%M:%S')
		
	except Exception, e:
		print exit_with_general_critical(e)

def index_miss_ratio(data):
	try:
		dbStatus['mongo.index.miss_ratio'] = float(data['indexCounters']['btree']['missRatio'])

	except Exception, e:
		print exit_with_general_critical(e)

def check_recordStatus(data):
	try:
		dbStatus['mongo.access.notinmemory'] = int(data['recordStats']['accessesNotInMemory'])
		dbStatus['mongo.page.faultexception'] = int(data['recordStats']['pageFaultExceptionsThrown'])
	except Exception, e:
		print exit_with_general_critical(e)

def check_replset_state(con):
	try:
		if not check_ismaster(con):
			set_read_preference(con.admin)
		data = con.admin.command(pymongo.son_manipulator.SON([('replSetGetStatus', 1)]))
	except:
		data = con.admin.command(son.SON([('replSetGetStatus', 1)]))
		
	primary_node = None
	secondary_node = []
	arbiter_node = []
	dbStatus['mongo.replication.replication_deplay'] = 0
	dbStatus['mongo.replication.replication_cluster_state'] = 1
	
	for member in data['members']:
		if member['stateStr'] == 'PRIMARY':
			primary_node = member
		if member['stateStr'] == 'SECONDARY':
			secondary_node.append(member)
		if member['stateStr'] == 'ARBITER':
			arbiter_node.append(member)
	if 	primary_node is not None:
		for member in secondary_node:
			if convert_time(primary_node['optime'].as_datetime()) - convert_time(member['optime'].as_datetime()) >300:
				dbStatus['mongo.replication.replication_deplay'] = convert_time(primary_node['optime'].as_datetime()) - convert_time(member['optime'].as_datetime())
		
	dbStatus['mongo.replication.replication_cluster_state'] = int(data['myState'])
	
def convert_time(obtime):
	t0 = obtime.timetuple()
	return time.mktime(t0)

def check_databases(con):
	try:
		try:
			if not check_ismaster(con):
				set_read_preference(con.admin)
			data = con.admin.command(pymongo.son_manipulator.SON([('listDatabases', 1)]))
		except:
			data = con.admin.command(son.SON([('listDatabases', 1)]))

		dbStatus['mongo.databases.count'] = len(data['databases'])
	except Exception, e:
		print exit_with_general_critical(e)

def check_collections(con):
	try:
		try:
			if not check_ismaster(con):
				set_read_preference(con.admin)
			data = con.admin.command(pymongo.son_manipulator.SON([('listDatabases', 1)]))
		except:
			data = con.admin.command(son.SON([('listDatabases', 1)]))

		count = 0
		for db in data['databases']:
			dbname = db['name']
			count += len(con[dbname].collection_names())

		dbStatus['mongo.collections.count'] = count

	except Exception, e:
		print exit_with_general_critical(e)


def check_all_databases_size(con):
	try:
		if not check_ismaster(con):
			set_read_preference(con.admin)
		all_dbs_data = con.admin.command(pymongo.son_manipulator.SON([('listDatabases', 1)]))
	except:
		all_dbs_data = con.admin.command(son.SON([('listDatabases', 1)]))

	total_storage_size = 0
	for db in all_dbs_data['databases']:
		database = db['name']
		data = con[database].command('dbstats') 
		storage_size = data['storageSize'] / 1024.0 / 1024.0
		total_storage_size += storage_size
		
	dbStatus['mongo.databases.total_storage_size'] = total_storage_size
	
def check_queues(data):
	try:
		dbStatus['mongo.queue.total_queues'] = float(data['globalLock']['currentQueue']['total']) 
		dbStatus['mongo.queue.readers_queues'] = float(data['globalLock']['currentQueue']['readers']) 
		dbStatus['mongo.queue.writers_queues'] = float(data['globalLock']['currentQueue']['writers']) 		
	except Exception, e:
		print exit_with_general_critical(e)

def check_oplog(con):
	try:
		db = con.local
		ol = db.system.namespaces.find_one({"name":"local.oplog.rs"})
		if (db.system.namespaces.find_one({"name":"local.oplog.rs"}) != None) :
			oplog = "oplog.rs"
		else:
			ol = db.system.namespaces.find_one({"name":"local.oplog.$main"})
			if (db.system.namespaces.find_one({"name":"local.oplog.$main"}) != None) :
				oplog = "oplog.$main"
		
		if not check_ismaster(con):
			set_read_preference(db)
		data = db.command('collstats', oplog)

		ol_size = data['size']
		ol_storage_size = data['storageSize']
		dbStatus['mongo.oplog.ol_used_storage'] = int(float(ol_size)/ol_storage_size*100+1)
		ol = con.local[oplog]
		firstc = ol.find().sort("$natural",pymongo.ASCENDING).limit(1)[0]['ts']
		lastc = ol.find().sort("$natural",pymongo.DESCENDING).limit(1)[0]['ts']
		time_in_oplog= (lastc.as_datetime()-firstc.as_datetime())
		try:
			dbStatus['mongo.oplog.hours_in_oplog'] = time_in_oplog.total_seconds()/60/60
		except:
			dbStatus['mongo.oplog.hours_in_oplog'] = float(time_in_oplog.seconds + time_in_oplog.days * 24 * 3600)/60/60

	except Exception, e:
		print exit_with_general_critical(e)

def check_journal_commits_in_wl(data):
	try:
		dbStatus['mongo.dur.commitsInWriteLock'] = data['dur']['commitsInWriteLock'] 

	except Exception, e:
		print exit_with_general_critical(e)

def check_journaled(data):
	try:
		dbStatus['mongo.dur.journaled'] = data['dur']['journaledMB'] 

	except Exception, e:
		print exit_with_general_critical(e)

def check_write_to_datafiles(data):
	try:
		dbStatus['mongo.dur.writes'] = data['dur']['writeToDataFilesMB'] 

	except Exception, e:
		print exit_with_general_critical(e)


def get_opcounters(data, opcounters_name, host, port):
	try : 
		insert = data[opcounters_name]['insert']
		query = data[opcounters_name]['query']
		update = data[opcounters_name]['update']
		delete = data[opcounters_name]['delete']
		getmore = data[opcounters_name]['getmore']
		command = data[opcounters_name]['command']
	except KeyError,e:
		return 0, [0]*100
	total_commands = insert + query + update + delete + getmore + command
	new_vals = [total_commands, insert, query, update, delete, getmore, command]
	return  maintain_delta(new_vals, host, port, opcounters_name)
	
def check_opcounters(data, host, port):
	err1, delta_opcounters = get_opcounters(data, 'opcounters', host, port)
	if err1 == 0:
		#per_minute_delta = [int(x/delta_opcounters[0]*60) for x in delta_opcounters[1:]]
		ops = delta_opcounters[1:]
		dbStatus['mongo.opcounter.total_commands'] = ops[0]
		dbStatus['mongo.opcounter.insert'] = ops[1]
		dbStatus['mongo.opcounter.query'] = ops[2]
		dbStatus['mongo.opcounter.update'] = ops[3]
		dbStatus['mongo.opcounter.delete'] = ops[4]
		dbStatus['mongo.opcounter.getmore'] = ops[5]
		dbStatus['mongo.opcounter.command'] = ops[6]
	else :
		print exit_with_general_critical("problem reading data from temp file")

def check_replica_primary(data):
	dbStatus['mongo.replication.current_primary'] = data['repl'].get('primary')
	
def build_file_name(host, port, action):
	basedir = os.getenv("HOME") + '/data/zabbix/'
	module_name = re.match('(.*//*)*(.*)\..*',__file__).group(2)
	return basedir + module_name + "-" + host + "-" + str(port) +"-" + action + ".data"

def ensure_dir(f):
	d = os.path.dirname(f) 
	if not os.path.exists(d):
		os.makedirs(d)
    
def write_values(file_name, string):
	f = None
	try:
		f = open(file_name, 'w')
	except IOError,e:
		if (e.errno == 2):
			ensure_dir(file_name)
			f = open(file_name, 'w')
		else:
			raise IOError(e)
	f.write(string)
	f.close()
	return 0
    
def read_values(file_name):
	data = None
	try:
		f = open(file_name, 'r')
		data = f.read()
		f.close()
		return 0, data
	except IOError, e:
		if (e.errno == 2):
			return 1, ''
	except Exception, e:
		return 2, None

def calc_delta(old, new):
	"""求差"""
	delta = []
	if (len(old) != len(new)):
		raise Exception("unequal number of parameters")
	for i in range(0,len(old)):
		val = float(new[i]) - float(old[i])
		if val < 0:
			val = new[i]      
		delta.append(val) 
	return 0, delta

def maintain_delta(new_vals, host, port, action):
	file_name = build_file_name(host, port, action)
	err, data = read_values(file_name)
	old_vals = data.split(';')
	new_vals = [str(int(time.time()))] + new_vals
	delta = None
	try:
		err, delta = calc_delta(old_vals, new_vals)
	except:
		err = 2
	write_res = write_values(file_name, ";".join(str(x) for x in new_vals))
	return err + write_res, delta
	
def get_mongod_data(data, con ,host, port):
	try:
		check_connections(data)
	except Exception, e:
		print e
	try:
		check_replset_state(con)
	except Exception, e:
		print e
	try:
		check_memory(data)
	except Exception, e:
		print e
	try:
		check_queues(data)
	except Exception, e:
		print e
	try:
		check_lock(data)
	except Exception, e:
		print e
	try:
		check_flushing(data)
	except Exception, e:
		print e
	try:
		index_miss_ratio(data)
	except Exception, e:
		print e
	try:
		check_databases(con)
	except Exception, e:
		print e
	try:
		check_collections(con)
	except Exception, e:
		print e
	try:
		check_oplog(con)
	except Exception, e:
		print e
	try:
		check_journal_commits_in_wl(data)
	except Exception, e:
		print e
	try:
		check_all_databases_size(con)
	except Exception, e:
		print e
	try:
		check_journaled(data)
	except Exception, e:
		print e
	try:
		check_write_to_datafiles(data)
	except Exception, e:
		print e
	try:
		check_opcounters(data, host, port)
	except Exception, e:
		print e
	try:
		check_replica_primary(data)
	except Exception, e:
		print e
	try:
		check_recordStatus(data)
	except Exception, e:
		print e
		
def get_mongos_data(data, con ,host, port):
	try:
		check_connections(data)
	except Exception, e:
		print e
	try:
		check_collections(con)
	except Exception, e:
		print e
	try:
		check_databases(con)
	except Exception, e:
		print e
	try:
		check_all_databases_size(con)
	except Exception, e:
		print e
	try:
		check_memory(data)
	except Exception, e:
		print e
	try:
		check_opcounters(data, host, port)
	except Exception, e:
		print e
		
def get_data(data, con ,host, port):
	if data['process'] == 'mongos':
		get_mongos_data(data, con ,host, port)
	elif data['process'] == 'mongod':
		get_mongod_data(data, con ,host, port)
	else:
		print 'process is not mongos or mongod'

def zabbixSender(argv):
	COMMAND_NAME = os.path.basename(sys.argv[0])
	p = optparse.OptionParser()
	
	p.add_option('-H', '--host', action='store', type='string', dest='host', default='127.0.0.1', help='The hostname you want to connect to')
	p.add_option('-P', '--port', action='store', type='int', dest='port', default=27017, help='The port mongodb is runnung on')
	p.add_option('-A', '--action', action='store', type='choice', dest='action', default='connections', help='The action you want to take',
					choices=['connections', 'replset_state', 'memory', 'lock', 'flushing', 'index_miss_ratio', 
							'databases', 'collections', 'database_size','queues','oplog','journal_commits_in_wl',
							'write_data_files','journaled','opcounters','replica_primary'])
	options, arguments = p.parse_args()
	
	host = options.host
	port = options.port		
	action = options.action
	
	DATA_FILE = '%s/data/zabbix/%s_%s_%s.data' % (os.getenv("HOME"), COMMAND_NAME, host, port)
	LOG_FILE = '%s/data/zabbix/%s_%s_%s.log' % (os.getenv("HOME"), COMMAND_NAME, host, port)
	
	with open(DATA_FILE, 'w') as fd:
		fd.truncate()
		fd.close()
	err, con = mongo_connect(host, port)
	
	if err != 0:
		return err
		
	data = get_server_status(con)
	get_data(data, con, host, port)
	with open(DATA_FILE, 'a') as fd:
		for k, v in dbStatus.items():
			line = "%s %s %s\n" % (host, k, v)
			fd.write(line)
	fd.close()
	cmd = '%s/bin/zabbix_sender -vv -c %s  -i %s >>%s 2>&1' % (os.getenv("HOME"), zabbix_agentd_config, DATA_FILE, LOG_FILE)
	os.system(cmd)
	
if __name__ == "__main__":
	sys.exit(zabbixSender(sys.argv[1:]))
