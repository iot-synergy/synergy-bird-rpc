echo "Build start " `date`
git config --global url.git@github.com:.insteadOf https://github.com

export GO111MODULE=on
export GOPRIVATE="github.com/iot-synergy"
rm -f mms-api
echo "start go mod tidy" `date`
go mod tidy --compat=1.21 
echo "start go build" `date`
go build  -o mms-api mms.go
echo "Build project for Linux successfully" `date`