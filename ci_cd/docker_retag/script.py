import os, docker, json
acr_endpoint = os.environ['INPUT_ACR_ENDPOINT']
image_old_repo = os.environ['INPUT_DOCKER_OLD_REPO']
image_new_repo = os.environ['INPUT_DOCKER_NEW_REPO']
new_tag = os.environ['INPUT_DOCKER_NEW_TAG']
old_tag = os.environ['INPUT_DOCKER_OLD_TAG']

#build_args = os.environ['INPUT_DOCKERBUILD_ARGS']

try:
    azure_credentials = os.environ['INPUT_AZURE_CREDENTIALS']
except:
    azure_credentials = 'NULL'

try:
    azure_username = os.environ['INPUT_AZURE_USERNAME']
    azure_password = os.environ['INPUT_AZURE_PASSWORD']
except:
    azure_username = 'NULL'
    azure_password = 'NULL'
    
try:
    clientDocker = docker.from_env()
except:
    print("There was an error using servers docker. Please check if docker is running.")

image_old_repo_path = acr_endpoint + '/' + image_old_repo
image_new_repo_path = acr_endpoint + '/' + image_new_repo


if azure_credentials != "NULL" :
    credentials = json.loads(azure_credentials)
    clientDocker.login(
        username=credentials["clientId"],
        password=credentials["clientSecret"],
        registry=acr_endpoint
    )
else:
    clientDocker.login(
            username=azure_username,
            password=azure_password,
            registry=acr_endpoint
    )

clientDocker.images.pull(
        repository=image_old_repo_path,
        tag=old_tag
)

image = clientDocker.images.get(image_old_repo_path + ":" + old_tag)
image.tag(image_new_repo_path,new_tag)

clientDocker.images.push(
        repository=image_new_repo_path,
        tag=new_tag
)