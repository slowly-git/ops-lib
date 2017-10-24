import subprocess

try:
    out_bytes = subprocess.check_output(['netstat', '-lpnt'], stderr=subprocess.STDOUT)
except subprocess.CalledProcessError as e:
    out_bytes = e.output

print(out_bytes.decode('utf-8'))

# 如果想让命令以shell来运行

try:
    shell_out_bytes = subprocess.check_output('grep python | wc > out', shell=True)
except subprocess.CalledProcessError as e:
    shell_out_bytes = e.output

print(shell_out_bytes.decode('utf-8'))

# 如果想自己控制subprocess的标准输入，需要用更底层的方式来实现

text = b"'hello world this is a test goodby'"
p = subprocess.Popen(['wc'],stdout=subprocess.PIPE, stdin=subprocess.PIPE)
stdout, stderr = p.communicate(text)
out = stdout.decode('utf-8')
err = stderr.decode('utf-8')

print(out, err)
