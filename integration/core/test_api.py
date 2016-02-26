import pytest
import cattle


@pytest.fixture
def client():
    url = 'http://localhost:8899/v2/schemas'
    return cattle.from_env(url=url)

def test_service_list(client):
    services = client.list_service()
    assert len(services) > 0

    assert services[0].name is not None

def test_container_create(client):
    restart_policy = {"maximumRetryCount": 2, "name": "on-failure"}
    health_check = {"name": "check1", "responseTimeout": 3,
                    "interval": 4, "healthyThreshold": 5,
                    "unhealthyThreshold": 6, "requestLine": "index.html",
                    "port": 200}
    build={
        'dockerfile': 'test/Dockerfile',
        'remote': 'http://example.com',
        'rm': True,
    }
    c = client.create_container(image="docker:ubuntu:latest",
     name="foo",
     restartPolicy=restart_policy,
     tty=True,
     startOnCreate=True,
     healthCheck=health_check,
     requestedIpAddress="10.1.1.19",
     privileged=True,
     domainName="rancher.io",
     #memory=8000000,
     stdinOpen=True,
     entryPoint=["/bin/sh", "-c"],
     cpuShares=400,
     cpuSet="0",
     workDir="/",
     hostname="test",
     user="test",
     environment={'TEST_FILE': "/etc/testpath.conf"},
     command=['sleep', '42'],
     capAdd=["SYS_MODULE"],
     capDrop=["SYS_MODULE"],
     build=build)

    c = client.wait_success(c)
    
    assert c.image == "docker:ubuntu:latest"
    assert c.name == "foo"
    assert c.restartPolicy is not None
    assert c.restartPolicy.name == 'on-failure'
    assert c.restartPolicy.maximumRetryCount == 2
    assert c.tty == True
    assert c.state != ""
    assert c.healthCheck is not None
    assert c.healthCheck.name == "check1" 
    assert c.healthCheck.responseTimeout == 3
    assert c.healthCheck.interval == 4
    assert c.healthCheck.healthyThreshold == 5
    assert c.healthCheck.unhealthyThreshold == 6
    assert c.healthCheck.requestLine == "index.html"
    assert c.healthCheck.port == 200
    assert c.startCount >=0
    assert c.revision == "0"
    assert c.startOnCreate == True
    #assert c.ipAddress != ""
    #assert c.firstRunning != ""
    assert c.nativeContainer == False
    assert c.token != ""
    assert c.externalId != ""
    #assert c.requestedIpAddress != ""
    assert c.privileged is True
    assert c.domainName == "rancher.io"
    #assert c.memory == 8000000
    assert c.stdinOpen is True
    assert c.entryPoint == ["/bin/sh", "-c"]
    assert c.cpuShares == 400
    assert c.cpuSet == "0"
    assert c.workDir == "/"
    assert c.hostname == "test"
    assert c.user == "test"
    assert c.environment == {'TEST_FILE': "/etc/testpath.conf"}
    #assert c.command == ['sleep', '42']
    assert c.publishAllPorts is not None
    assert c.capAdd == ["SYS_MODULE"]
    assert c.capDrop == ["SYS_MODULE"]
    assert c.build.dockerfile == 'test/Dockerfile'

