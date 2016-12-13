package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

)

func checkIP() (string, error) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(buf)), nil
}

var(
	sgid string
	region string
	wait string
	port int64
)

func main() {
	ip, err := checkIP()

	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Enter Region: ")
	fmt.Scanln(&region)
	fmt.Printf("Enter Security Group-Id: ")
	fmt.Scanln(&sgid)
	fmt.Printf("Door on: (port)")
	fmt.Scanln(&port)
	fmt.Printf("My IP is %q\n", ip)
	fmt.Printf("Open Sesame\n")
	
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
	
	// Add security group rule
	authsgin := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    aws.String(sgid),
		CidrIp:     aws.String(ip + "/32"),
		DryRun:     aws.Bool(false),
		FromPort:   aws.Int64(port),
		ToPort:     aws.Int64(port),
		IpProtocol: aws.String("tcp"),
		}
	resp, err := svc.AuthorizeSecurityGroupIngress(authsgin)

	if err != nil {
	// Print the error, cast err to awserr.Error to get the Code and
	// Message from an error.
	fmt.Println(err.Error())
	return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	
	// Waiting for operation finished
	// Scanln(wait) for keep the firewall open
	fmt.Printf("Waiting for the door close...(Enter anything if operation finished)")
	fmt.Scanln(&wait)
	
	// Revoke security group rule
	revsgin := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:    aws.String(sgid),
		CidrIp:     aws.String(ip + "/32"),
		DryRun:     aws.Bool(false),
		FromPort:   aws.Int64(port),
		ToPort:     aws.Int64(port),
		IpProtocol: aws.String("tcp"),
		}
	svc.RevokeSecurityGroupIngress(revsgin)
	
	fmt.Printf("Close Sesame")
}
