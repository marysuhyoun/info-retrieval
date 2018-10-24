
//This go file creates indexes of the collection with the option of stemming.

package main

import (
"os"
"bufio"
s "strings"
"strconv"
"bytes"
"encoding/binary"
	"fmt"
	"time"
)
//hash map for each dictionary item (tuple : postingMap)
var invertedIndex map[string]map[int]int
var binaryOutput string
var dictionary map[string][]int
var outputFileName string ="/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14StemInvertedIndex.txt"
var outputFile,_ = os.Create(outputFileName)
var totalDocs int
var totalWords int
var stemming string
var option string

func main() {
	fmt.Println("Inverted Index Creation Start Time: ", time.Now())
	option = "writeFiles"
	//Scan input file and populate invertedIndex
	file := "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14.txt"
	stemming = "yes"
	InitializeInvertedIndex()
	//generate dictionary of word : offset , length
	InitializeDictionary()
	totalWords = 0
	ScanFile(file)
	//populate dictionary and append inverted index
	offset := 0
	i:= 0
	for key, value := range invertedIndex{
		offset = UpdateDictionary(key, value, offset)
		i++
	}
	if option == "writeFiles" {
		//create binary file output to write into
		outputFile.Close()


		fmt.Println("Total Number of Documents" , totalDocs)
		fmt.Println("done")

		//write dictionary into an output file
		WriteDictionary(dictionary)
	}
	fmt.Println("Inverted Index Creation End Time: ", time.Now())
}

//initialize dictionary
func InitializeDictionary(){
	dictionary = make(map[string][]int)
	return
}

//scan the file line by line
func ScanFile(filePath string,){
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var currentID int
	for scanner.Scan(){
		nextLine := scanner.Text()
		// </P> marks the end of the document.
		if !s.Contains(nextLine, "</P>"){
			//determine docID
			if s.HasPrefix(nextLine, "<P ID="){
				//remove prefix and suffix of docID
				nextLine = s.TrimPrefix(nextLine, "<P ID=")
				nextLine = s.TrimSuffix(nextLine, ">")
				currentID,_ = strconv.Atoi(nextLine)
				totalDocs = totalDocs + 1
			}else if nextLine != "" { //process if current line is not blank. otherwise skip this line.
				//normalize text (non-stemming)
				normalizedLine := NormalizeText(nextLine)
				//remove stop words from the list. https://gist.github.com/sebleier/554280
				normalizedLine = RemoveStopWords(normalizedLine)
				//split up current line into an array of strings
				lineArray := s.Split(normalizedLine, " ")
				//stem if required
				if stemming == "yes" {
					StemText(lineArray)
				}
				//update inverted index
				UpdateIndex(lineArray,currentID)
			}
		}
	}
}
//given set of text, normalize it
func NormalizeText(lineContent string) string{
	lineContent = s.Replace(lineContent,".", "",-1)
	lineContent = s.Replace(lineContent,"-", " ",-1)
	lineContent = s.Replace(lineContent,",", "",-1)
	lineContent = s.Replace(lineContent,"?", "",-1)
	lineContent = s.Replace(lineContent,";", "",-1)
	lineContent = s.Replace(lineContent,":", "",-1)
	lineContent = s.Replace(lineContent,"(", "",-1)
	lineContent = s.Replace(lineContent,")", "",-1)
	lineContent = s.Replace(lineContent,"[", "",-1)
	lineContent = s.Replace(lineContent,"]", "",-1)
	lineContent = s.Replace(lineContent,"{", "",-1)
	lineContent = s.Replace(lineContent,"}", "",-1)
	lineContent = s.Replace(lineContent,"=", "",-1)
	lineContent = s.Replace(lineContent,"_", "",-1)
	lineContent = s.Replace(lineContent,"+", "",-1)
	lineContent = s.Replace(lineContent,"$", "",-1)
	lineContent = s.Replace(lineContent,"#", "",-1)
	lineContent = s.Replace(lineContent,"@", " ",1 )
	lineContent = s.Replace(lineContent,"|", "",-1)
	lineContent = s.Replace(lineContent,"\\n", "",-1)
	lineContent = s.Replace(lineContent,"- ", " ",-1)
	lineContent = s.Replace(lineContent,"'s", "",-1)
	lineContent = s.Replace(lineContent,"'ve", "",-1)
	lineContent = s.Replace(lineContent,"'m", "",-1)
	lineContent = s.Replace(lineContent,"'re", "",-1)
	lineContent = s.Replace(lineContent,"'ed", "",-1)
	lineContent = s.Replace(lineContent,"'", "",-1)
	lineContent = s.Replace(lineContent,"  ", " ",-1)
	lineContent = s.Replace(lineContent,"  ", " ",-1)
	lineContent = s.Replace(lineContent,"!", "",-1)
	lineContent = s.Replace(lineContent,":", "",-1)
	lineContent = s.Replace(lineContent,"...", "",-1)
	lineContent = s.Replace(lineContent,"<", "",-1)
	lineContent = s.Replace(lineContent,">", "",-1)
	lineContent = s.Replace(lineContent,"/", "",-1)
	lineContent = s.Replace(lineContent,"*", "",-1)
	lineContent = s.Replace(lineContent,"%", "",-1)
	lineContent = s.Replace(lineContent,"^", "",-1)
	lineContent = s.Replace(lineContent,"_", "",-1)

	return s.ToLower(lineContent)
}

//remove stop words
func RemoveStopWords(lineContent string) string {
	stopwords := [127]string{"i","me","my","myself","we","our","ours","ourselves","you","your","yours","yourself","yourselves","he","him","his","himself","she","her","hers","herself","it","its","itself","they","them","their","theirs","themselves","what","which","who","whom","this","that","these","those","am","is","are","was","were","be","been","being","have","has","had","having","do","does","did","doing","a","an","the","and","but","if","or","because","as","until","while","of","at","by","for","with","about","against","between","into","through","during","before","after","above","below","to","from","up","down","in","out","on","off","over","under","again","further","then","once","here","there","when","where","why","how","all","any","both","each","few","more","most","other","some","such","no","nor","not","only","own","same","so","than","too","very","s","t","can","will","just","don","should","now"}
	for i := 0 ; i < len(stopwords) ; i++{
		if s.Contains(lineContent, stopwords[i]){
			s.Replace(lineContent, stopwords[i],"", -1)
		}
	}
	return lineContent
}
//given set of text, stem it to the first 5 characters
func StemText(lineArray []string) []string{
	for i:=0 ; i < len(lineArray) ; i ++{
		if len(lineArray[i]) > 5 {
			var output rune
			for pos, char := range lineArray[i] {
				if pos < 5 {
					output = output + char
				}
			}
			lineArray[i] = strconv.QuoteRune(output)
		}
	}
	return lineArray
}
//initialize invertedIndex
func InitializeInvertedIndex(){
	invertedIndex = make(map[string]map[int]int)
	return
}

//update the dictionary with input word and docID
func UpdateIndex(lineArray []string , docID int){
	//make a postingsMap for every unique word in the lineArray and add it to the invertedIndex
	for i := 0 ; i < len(lineArray) ; i++ {
		//make sure the string in question is not composed of just " " spaces.
		lookUpItem := s.TrimSpace(lineArray[i])
		if lookUpItem != "" {
			//see if tuple already exists in invertedIndex
			if key, ok := invertedIndex[lineArray[i]]; ok{
				//see if docID already exists in keyMap (postingsMap)
				if _, ok := key[docID]; ok {
					key[docID] = key[docID] + 1
					invertedIndex[lineArray[i]] = key
				}else { //create new item if docID does not exist in keyMap (postingsMap)
					key[docID] = 1
					invertedIndex[lineArray[i]] = key
				}
			} else { //create a new tuple item if one does not exist
				postingsMap := make(map[int]int)
				postingsMap[docID] = 1
				invertedIndex[lineArray[i]] = make(map[int]int)
				invertedIndex[lineArray[i]] = postingsMap
			}
		}

	}
}
func UpdateDictionary(term string, postingsMap map[int]int, offset int) int{
	lenPostingsMap := len(postingsMap)
	dictionary[term] = []int{lenPostingsMap,offset}
	newOffset := offset
	for key, value := range postingsMap {
		//get key values in binary
		keyBinary := make([]byte, 4)
		binary.LittleEndian.PutUint32(keyBinary, uint32(key))

		valueBinary := make([]byte, 4)
		binary.LittleEndian.PutUint32(valueBinary, uint32(value))
		if option == "writeFiles"{
			WriteToFile(keyBinary)
			WriteToFile(valueBinary)
		}
		newOffset = newOffset + 8
	}
	return newOffset
}
func WriteToFile (input []byte) {

	f, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, input[0])
	binary.Write(buf, binary.LittleEndian, input[1])
	binary.Write(buf, binary.LittleEndian, input[2])
	binary.Write(buf, binary.LittleEndian, input[3])
	output := make([]byte, 4)
	output[0] = input[0]
	output[1] = input[1]
	output[2] = input[2]
	output[3] = input[3]
	if _, err = f.Write(output); err != nil {
		panic(err)
	}
}
func WriteDictionary (dictionary map[string][]int){
	var dictionaryFileName string ="/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14StemDictionary.txt"
	dictionaryFile, _ := os.Create(dictionaryFileName)
	dictionaryFile.Close()
	f, err := os.OpenFile(dictionaryFileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	if _, err = f.WriteString("key , " + "length" + " , " + " offset" + "\n"); err != nil {
		panic(err)
	}
	var stringToPrint string
	for key, value := range dictionary {
		stringToPrint = key + " , " + strconv.Itoa(value[0]) + " , " + strconv.Itoa(value[1]) + "\n"
		if _, err = f.WriteString(stringToPrint); err != nil {
			panic(err)
		}
	}
	dictionaryFile.Close()


}