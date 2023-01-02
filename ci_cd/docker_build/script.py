import os, docker, json
acr_endpoint = os.environ['INPUT_ACR_ENDPOINT']
azure_credentials = os.environ['INPUT_AZURE_CREDENTIALS']
build_path = os.environ['INPUT_DOCKERBUILD_PATH']
build_repo = os.environ['INPUT_DOCKERBUILD_REPO']
build_tag = os.environ['INPUT_DOCKERBUILD_TAG']
build_args = os.environ['INPUT_DOCKERBUILD_ARGS']

credentials = json.loads(azure_credentials)
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

print(credentials["clientId"])
print(credentials["clientSecret"])
print(acr_endpoint)
try:
    clientDocker.login(
        username=credentials["clientId"],
        password=credentials["clientSecret"],
        registry=acr_endpoint
    )
except:
    print(credentials)

# try:
#     pushResponse = clientDocker.api.push(
#         repository=build_repo_path,
#         tag=build_tag,
#     )
# except:
#     print(credentials)
