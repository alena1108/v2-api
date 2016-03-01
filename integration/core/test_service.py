import pytest
import cattle
from random import choice
from string import ascii_uppercase

@pytest.fixture
def client():
    url = 'http://localhost:8899/v2/schemas'
    return cattle.from_env(url=url)

def get_stack_id():
    return "1s1"

def test_service_create_basic(client):
    metadata = {"bar": {"people": [{"id": 0}]}}
    name = ''.join(choice(ascii_uppercase) for i in range(4));
    l_sel = "foo=bar"
    c_sel = "bar=foo"
    s = client.create_service(name=name,
        stackId= get_stack_id(),
        scale=2,
        serviceIpAddress='10.1.1.1',
        assignServiceIpAddress=True,
        metadata=metadata,
        linkSelector=l_sel,
        containerSelector=c_sel,
        retainIpAddress=True)

    s = client.wait_success(s)
    
    assert s.state == 'inactive'
    assert s.stackId ==  get_stack_id()
    assert s.name == name
    assert s.createIndex == 0
    assert s.scale == 2
    assert s.metadata == metadata
    assert s.linkSelector == l_sel
    assert s.containerSelector == c_sel
    assert s.retainIpAddress == True


def test_service_lc(client):
    name = ''.join(choice(ascii_uppercase) for i in range(4));
    lc = {"image": "docker:ubuntu:latest"}
    s = client.create_service(name=name,
        stackId= get_stack_id(),
        containerTemplates=[lc])

    s = client.wait_success(s)
    
    assert s.state == 'inactive'
    assert s.stackId ==  get_stack_id()
    assert s.name == name
    assert s.containerTemplates == [lc]
