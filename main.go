//** Untested beyone 1000 keys to delete, which is the limit of the bulk delete function

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type flagArray map[string]bool

func (i flagArray) Set(value string) error {
	i[value] = true
	return nil
}

//Needed to match the flag.Var interface
func (i flagArray) String() string {
	return ""
}

func main() {
	fmt.Println("Started execution...")
	bucketName := flag.String("bucket", "", "Bucket Name to bust cache of")
	fileName := flag.String("file", "", "file whose cache needs busting")
	dirToIgnore := make(flagArray)
	flag.Var(&dirToIgnore, "ignore", "directories to ignore≤, omit trailing and leading slashes")
	flag.Parse()

	sess, errSession := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if errSession != nil {
		fmt.Printf("Error creating session: %s", errSession)
		os.Exit(1)
	}

	client := s3.New(sess)

	keysToDelete, pageNum, totalRecords := GetKeys(client, bucketName, fileName, dirToIgnore)

	if len(keysToDelete) == 0 {
		fmt.Println("No keys found to delete")
		os.Exit(0)
	}

	buf := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to proceed (y/n): ")
	input, err := buf.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}

	if !(len(input) == 2 && input[0] == 'y') {
		fmt.Println("Stopping execution, nothing will be deleted")
		os.Exit(0)
	}

	DeleteKeys(client, bucketName, keysToDelete)

	fmt.Printf("Done: %d pages, %d records\n", pageNum, totalRecords)

}

func GetKeys(client *s3.S3, bucketName *string, fileName *string, ignore flagArray) ([]*s3.ObjectIdentifier, int, int) {

	bucketRequest := &s3.ListObjectsV2Input{
		Bucket:  aws.String(*bucketName),
		MaxKeys: aws.Int64(100),
	}

	pageNum := 0
	totalRecords := 0
	var keysToDelete = []*s3.ObjectIdentifier{}

	errRead := client.ListObjectsV2Pages(bucketRequest,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			for _, e := range page.Contents {
				totalRecords++
				if e.Key == nil {
					break
				}
				var key string
				key = *e.Key
				keyFileName := filepath.Base(key)
				dir := filepath.Dir(key)

				if ignore[dir] {
					if keyFileName == *fileName {
						fmt.Print("Don't delete ignored records: ")
						fmt.Println(key)
					}
				} else if keyFileName == *fileName {
					keyToDelete := s3.ObjectIdentifier{
						Key: e.Key,
					}
					keysToDelete = append(keysToDelete, &keyToDelete)
					fmt.Print("Queued to delete: ")
					fmt.Println(key)
				}
			}
			return !lastPage
		})

	if errRead != nil {
		fmt.Println(errRead)
		os.Exit(2)
	}
	return keysToDelete, pageNum, totalRecords
}

func DeleteKeys(client *s3.S3, bucketName *string, keysToDelete []*s3.ObjectIdentifier) {
	deleteKeys := s3.Delete{
		Objects: keysToDelete,
	}

	deleteRequestInput := s3.DeleteObjectsInput{
		Bucket: aws.String(*bucketName),
		Delete: &deleteKeys,
	}

	deleteReq, deleteResp := client.DeleteObjectsRequest(&deleteRequestInput)

	deleteErr := deleteReq.Send()
	if deleteErr != nil {
		fmt.Println(deleteErr)
		os.Exit(3)
	}

	fmt.Println(deleteResp)
}
