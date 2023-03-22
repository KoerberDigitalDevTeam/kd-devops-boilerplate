import os, docker, json, sys
from subprocess import run, PIPE

acr_endpoint = os.environ['INPUT_ACR_ENDPOINT']
build_path = os.environ['INPUT_DOCKERBUILD_PATH']
build_repo = os.environ['INPUT_DOCKERBUILD_REPO']
build_tag = os.environ['INPUT_DOCKERBUILD_TAG']
build_args = os.environ['INPUT_DOCKERBUILD_ARGS']

try:
    docker_scan = os.environ['INPUT_DOCKER_SCAN']
except:
    docker_scan = "false"

try:
    docker_trivy_image_flags = os.environ['INPUT_DOCKER_TRIVY_IMAGE_FLAGS']
    docker_trivy_vulnerability_ignore = os.environ['INPUT_DOCKER_TRIVY_VULNERABILITY_IGNORE']
except:
    docker_trivy_image_flags = "--exit-code 1 --severity HIGH,CRITICAL --ignore-unfixed"
    docker_trivy_vulnerability_ignore = ""

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

build_repo_path = acr_endpoint + '/' + build_repo
build_complete_path = build_repo_path + ':' + build_tag
build_complete_path = build_complete_path.lower()
clientDocker.images.build(
        path=build_path,
        tag=build_complete_path,
        buildargs=json.loads(build_args)
)

if docker_scan == "true":
    vulnerability_ignored = docker_trivy_vulnerability_ignore.split(",")
    ignore_vun = open(".trivyignore","w")
    for vulnerability in vulnerability_ignored:
        ignore_vun.write(vulnerability + "\n")
    ignore_vun.close()

    # exit_code_os = os.system('trivy image' + ' ' + docker_trivy_image_flags + ' ' + build_complete_path)
    # exit_code = exit_code_os >> 8
    # if exit_code != 00000000: 
    #     sys.exit('Something bad happened')

    p = run( ['trivy', 'image ' + docker_trivy_image_flags + ' ' +build_complete_path ], stdout=PIPE, stderr=PIPE, text=True )
    if p.returncode != 0:
        print('Something bad happened')
    print("Output")
    print(p.stdout)
    print("It is work")
    print(p.stderr)

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
