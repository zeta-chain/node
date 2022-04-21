ECR_REGISTRY=821351686724.dkr.ecr.us-east-1.amazonaws.com
ECR_REPOSITORY=zetachain-sparta
IMAGE_TAG=zetacore

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 821351686724.dkr.ecr.us-east-1.amazonaws.com

docker build -f Dockerfile.zetacore -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG