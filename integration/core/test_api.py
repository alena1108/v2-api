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
    c = client.create_container(image="docker:ubuntu:latest",
     name="tt",
     restartPolicy=restart_policy,
     startOnCreate=True,
     tty=True,
     healthCheck=health_check)

    c = client.wait_success(c)
    
    assert c.image == "docker:ubuntu:latest"
    assert c.name == "tt"
    assert c.restartPolicy is not None
    assert c.restartPolicy.name == 'on-failure'
    assert c.restartPolicy.maximumRetryCount == 2
    assert c.tty == True
    assert c.state == "running"
    assert c.healthCheck is not None
    assert c.healthCheck.name == "check1" 
    assert c.healthCheck.responseTimeout == 3
    assert c.healthCheck.interval == 4
    assert c.healthCheck.healthyThreshold == 5
    assert c.healthCheck.unhealthyThreshold == 6
    assert c.healthCheck.requestLine == "index.html"
    assert c.healthCheck.port == 200
    assert c.startCount == 0
    assert c.revision == "0"
    assert c.startOnCreate == True
    assert c.ipAddress != ""
    assert c.firstRunning != ""
    assert c.nativeContainer == False

