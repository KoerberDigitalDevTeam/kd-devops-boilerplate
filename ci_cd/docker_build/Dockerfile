FROM python:3-alpine

# RUN apt-get update
# RUN apt-get install -y wget apt-transport-https gnupg lsb-release
# RUN wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor | tee /usr/share/keyrings/trivy.gpg > /dev/null
# RUN echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | tee -a /etc/apt/sources.list.d/trivy.list
# RUN apt-get update
# RUN apt-get install trivy

RUN apk --no-cache add curl
#Trivy
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.38.3
#Dockle
RUN arch=$(arch | sed s/aarch64/ARM64/ | sed s/x86_64/64bit/) && \
  wget "https://github.com/goodwithtech/dockle/releases/download/v0.4.11/dockle_0.4.11_Linux-${arch}.apk" && \
  apk add --allow-untrusted --no-cache dockle_0.4.11_Linux-${arch}.apk && \
  rm -rf dockle_0.4.11_Linux-${arch}.apk && \
  rm -rf /tmp/*


COPY requirements.txt /requirements.txt
RUN pip install -r /requirements.txt
COPY script.py /script.py
CMD ["python", "/script.py"]