import requests 
import json


def connect_peers():

    miner_info = requests.post("http://127.0.0.1:9545", json={"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}).json()['result']
    shutdown_info = requests.post("http://127.0.0.1:9547", json={"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}).json()['result']

    requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addPeer","params":[f"{miner_info['enode'].split('?')[0]}"],"id":1})
    requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addPeer","params":[f"{shutdown_info['enode'].split('?')[0]}"],"id":1})

if __name__ == "__main__":
    connect_peers()