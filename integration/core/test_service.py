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


def getName():
    return ''.join(choice(ascii_uppercase) for i in range(4))


def test_service_create_basic(client):
    metadata = {"bar": {"people": [{"id": 0}]}}
    l_sel = "foo=bar"
    c_sel = "bar=foo"
    name = getName()
    s = client.create_service(name=name,
                              stackId=get_stack_id(),
                              scale=2,
                              serviceIpAddress='10.1.1.1',
                              assignServiceIpAddress=True,
                              metadata=metadata,
                              linkSelector=l_sel,
                              containerSelector=c_sel,
                              retainIpAddress=True)

    s = client.wait_success(s)

    assert s.state == 'inactive'
    assert s.stackId == get_stack_id()
    assert s.name == name
    assert s.createIndex == 0
    assert s.scale == 2
    assert s.metadata == metadata
    assert s.linkSelector == l_sel
    assert s.containerSelector == c_sel
    assert s.retainIpAddress is True
    assert s.serviceIpAddress is not None


def test_service_lc(client):
    lc = {"image": "docker:ubuntu:latest"}
    name = getName()
    s = client.create_service(name=name,
                              stackId=get_stack_id(),
                              containerTemplates=[lc,lc])

    s = client.wait_success(s)

    assert s.state == 'inactive'
    assert s.stackId == get_stack_id()
    assert s.name == name
    assert s.containerTemplates is not None
    assert len(s.containerTemplates) == 2
    assert s.containerTemplates[0].image == "docker:ubuntu:latest"
    assert s.containerTemplates[1].image == "docker:ubuntu:latest"
