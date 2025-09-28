TAG=url-shortener
PUSH=$1


docker build -t $TAG . 
BUILD_RESULT=$?
if [ $BUILD_RESULT -ne 0 ]; then
    echo -e "\e[31mDocker build failed with exit code $BUILD_RESULT\e[0m"
    exit $BUILD_RESULT
fi
echo "Docker build completed successfully: $TAG"

if [ "$PUSH" = "push" ]; then
    echo "Pushing image $TAG"
    docker push $TAG
    exit 0
fi
