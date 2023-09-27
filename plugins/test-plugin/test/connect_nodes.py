import requests 
import json
import time


def connect_peers():

    miner_info = requests.post("http://127.0.0.1:9545", json={"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}).json()['result']
    tester_info = requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}).json()['result']
    shutdown_info = requests.post("http://127.0.0.1:9547", json={"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}).json()['result']

    miner_enode = f"{miner_info['enode'].split('?')[0].split('@')[0]}@127.0.0.1:64480"
    tester_enode = f"{miner_info['enode'].split('?')[0].split('@')[0]}@127.0.0.1:64481"
    shutdown_enode = f"{miner_info['enode'].split('?')[0].split('@')[0]}@127.0.0.1:64484"

    print(miner_enode)
    print(tester_enode)
    print(shutdown_enode)

    print(requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addTrustedPeer","params":[f"{miner_enode}"],"id":1}).text)
    print(requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addTrustedPeer","params":[f"{shutdown_enode}"],"id":1}).text)
    print(requests.post("http://127.0.0.1:9545", json={"jsonrpc":"2.0","method":"admin_addTrustedPeer","params":[f"{tester_enode}"],"id":1}).text)
    print(requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addPeer","params":[f"{miner_enode}"],"id":1}).text)
    print(requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"admin_addPeer","params":[f"{shutdown_enode}"],"id":1}).text)
    print(requests.post("http://127.0.0.1:9545", json={"jsonrpc":"2.0","method":"admin_addPeer","params":[f"{tester_enode}"],"id":1}).text)

    time.sleep(30)

    print(requests.post("http://127.0.0.1:9545", json={"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}).json()['result'])
    print(requests.post("http://127.0.0.1:9546", json={"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":1}).json()['result'])

if __name__ == "__main__":
    connect_peers()