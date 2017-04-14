## BUILD & RUN m-apiserver's image

```bash
# build & tag the image
sudo docker build -t openebs/m-apiserver:latest -t openebs/m-apiserver:0.2-RC4 .

# run the image as a docker container
sudo docker run -itd openebs/m-apiserver:latest

# verify the container
sudo docker ps
sudo docker ps -a
sudo docker logs <Container-ID>
sudo docker inspect <Container-ID>

# verify the m-apiserver service within this container
curl http://<Container-IP>:5656/latest/meta-data/instance-id

# run commands inside the container
amit:docker$ sudo docker exec -it 921f974ee490 bash
root@921f974ee490:/# 
root@921f974ee490:/# cat /etc/mayaserver/orchprovider/nomad_global.INI
```

## TODO

- Run m-apiserver up command without -bind option
- Follow dockerfile best practices
- Follow entrypoint best practices
- Follow scripting best practices
- Run linting on dockerfile
- Run linting on script
- Use a GoLang binary script than an entrypoint shell script
