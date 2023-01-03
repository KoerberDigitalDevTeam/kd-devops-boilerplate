import os, docker, json
acr_endpoint = os.environ['INPUT_ACR_ENDPOINT']
try:
    azure_credentials = os.environ['INPUT_AZURE_CREDENTIALS']
except:
    azure_credentials = 'NULL'
    
azure_username = os.environ['INPUT_AZURE_USERNAME']
azure_password = os.environ['INPUT_AZURE_PASSWORD']
build_path = os.environ['INPUT_DOCKERBUILD_PATH']
build_repo = os.environ['INPUT_DOCKERBUILD_REPO']
build_tag = os.environ['INPUT_DOCKERBUILD_TAG']
build_args = os.environ['INPUT_DOCKERBUILD_ARGS']

try:
    clientDocker = docker.from_env()
except:
    print("There was an error using servers docker. Please check if docker is running.")

build_repo_path = acr_endpoint + '/' + build_repo
build_complete_path = build_repo_path + ':' + build_tag
clientDocker.images.build(
        path=build_path,
        tag=build_complete_path,
        buildargs=json.loads(build_args)
)

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

pushResponse = clientDocker.api.push(
    repository=build_repo_path,
    tag=build_tag,
)
