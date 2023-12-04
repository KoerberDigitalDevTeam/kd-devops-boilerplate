import os, docker, json, sys
from subprocess import run, PIPE

acr_endpoint = os.environ['INPUT_ACR_ENDPOINT']
build_path = os.environ['INPUT_DOCKERBUILD_PATH']
build_dockerfile = os.environ['INPUT_DOCKERBUILD_FILE']

build_repo = os.environ['INPUT_DOCKERBUILD_REPO']
build_tag = os.environ['INPUT_DOCKERBUILD_TAG']
build_args = os.environ['INPUT_DOCKERBUILD_ARGS']

try:
    docker_scan = os.environ['INPUT_DOCKER_SCAN']
except:
    docker_scan = "false"

try:
    build_labels = os.environ['INPUT_DOCKERBUILD_LABELS']
except:
    build_labels = "{}"


try:
    docker_trivy_vulnerability_ignore = os.environ['INPUT_DOCKER_TRIVY_VULNERABILITY_IGNORE']
except:
    docker_trivy_vulnerability_ignore = ""

try:
    docker_trivy_image_flags = os.environ['INPUT_DOCKER_TRIVY_IMAGE_FLAGS']
except:
    docker_trivy_image_flags = "--exit-code 1 --severity HIGH,CRITICAL --ignore-unfixed"

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

print("Hello" + build_dockerfile)

build_repo_path = str(acr_endpoint + '/' + build_repo).lower()
build_complete_path = build_repo_path + ':' + str(build_tag).lower()
image, build_logs = clientDocker.images.build(
        path=build_path,
        dockerfile=build_dockerfile,
        tag=build_complete_path,
        labels=json.loads(build_labels),
        buildargs=json.loads(build_args),
        quiet=False
)
for chuck in build_logs:
    if 'stream' in chuck:
        for line in chuck['stream'].splitlines():
            print(line)

if docker_scan == "true":
    print("Image scan vulnerabilities")
    vulnerability_ignored = docker_trivy_vulnerability_ignore.split(",")
    ignore_vun = open(".trivyignore","w")
    for vulnerability in vulnerability_ignored:
        ignore_vun.write(vulnerability + "\n")
    ignore_vun.close()

    command = 'trivy image ' + docker_trivy_image_flags + ' ' + build_complete_path
    print("Exec -> " + command)
    p = run( command.split(), stdout=PIPE, stderr=PIPE, text=True )
    print(p.stdout)
    if p.returncode != 0:
        sys.exit('Something bad happened')

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