url=209479273800.dkr.ecr.ap-northeast-2.amazonaws.com/prod-cheego556-noti-ecr

aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-2.amazonaws.com"
docker pull $url:latest
docker stop noti || true
docker rm noti || true
docker run -d --name noti $url:latest
