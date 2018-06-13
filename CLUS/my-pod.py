#!/usr/bin/env python

import time

from kubernetes import config
from kubernetes.client import Configuration
from kubernetes.client.apis import core_v1_api
from kubernetes.client.rest import ApiException
from kubernetes.stream import stream

class My_pod(object):
	api = ""
	name = "" 
	conn = ""

	def __init__(self, name, host):
		self.api = core_v1_api.CoreV1Api()
		self.name = name + "-" + host
		resp = None
		try:
			resp = self.api.read_namespaced_pod(name=self.name,
							    namespace='default')
		except ApiException as e:
			if e.status != 404:
				print("Unknown error: %s" %e)
				exit(1)
		if not resp:
			print("Pod %s does not exist. Creating it..." % self.name)
			pod_manifest = {
 				'apiVersion': 'v1',
				'kind': 'Pod',
				'metadata': {
					'name': self.name
				},
			        'spec': {
					'containers': [{
						'image': name,
						'name': 'sleep',
						"args": [
							"/bin/sh",
							"-c",
							"while true;do date;sleep 5; done"
						]
					}],
					'restartPolicy': 'Always',
					'affinity': {
						'nodeAffinity': {
							'requiredDuringSchedulingIgnoredDuringExecution': {
								'nodeSelectorTerms': [{
									'matchExpressions': [{
										'key': 'dedicated',
										'operator': 'In',
										'values': [host]
									}]
								}]
							}

						}
					},
					"tolerations": [{
						'operator': 'Exists'
					}]
				}
			}

			resp = self.api.create_namespaced_pod(body=pod_manifest,
                                     			      namespace='default')
    			while True:
				resp = self.api.read_namespaced_pod(name=self.name,
								    namespace='default')
				if resp.status.phase != 'Pending':
					break
				time.sleep(1)
			print("Done.")
 
 	def open_conn(self):
		self.conn = stream(	self.api.connect_get_namespaced_pod_exec,
					self.name,
					'default',
					command=['/bin/sh'],
					stderr=True,
					stdin=True,
					stdout=True,
					tty=False,
					_preload_content=False)
	
 	def close_conn(self):
		self.conn.close()

 	def send(self, command):
		self.conn.update(timeout=1)
		self.conn.write_stdin(command + "\n")
		while self.conn.is_open():
			self.conn.update(timeout=1)
                	if self.conn.peek_stdout():
                        	print("%s" % self.conn.read_stdout())
				break

if __name__ == '__main__':
	pod = []
	config.load_kube_config()
	c = Configuration()
	c.assert_hostname = False
	Configuration.set_default(c)

	start_time = time.time()

	hosts = ['master', 'worker1']
	for i in range (0, len(hosts)):
		pod.append(My_pod('busybox', hosts[i]))
		pod[i].open_conn()

	start_time += time.time()

	'''
	# for k in range (2, 3):
        for k in range (2, 8):
		# for j in range (2, 101):
		for j in range (1, 255):
			start_time = time.time()
			for i in range(0, len(hosts)):
				# ip = "10.1." + str(k) + "." + str(j)
				ip = "10.96." + str(k) + "." + str(j)  
			pod[i].send("wget -O - " + ip + "/server_addr:80")
			pod[i].send("wget -O - " + ip + "/server_addr:8000")
			pod[i].send("wget -O - " + ip + "/server_addr:80")
			# pod[i].send("wget -O - 10.96.1.1/server_addr:8000")
			# pod[i].send("wget -O - 10.96.1.1/server_addr:8080")
			# pod[i].send("wget -O - 10.96.1.1/server_addr:8080")
			print time.time() - start_time
	'''

	for i in range(0, len(hosts)):
		pod[i].close_conn()
