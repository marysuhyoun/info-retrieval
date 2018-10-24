package Programming_Assignments

import (
	"os"
	"bufio"
	s "strings"
	"strconv"
	"fmt"
	"sort"
	"bytes"
	"encoding/binary"
)
//hash map for each dictionary item (tuple : postingMap)
var invertedIndex map[string]map[int]int
var binaryOutput string
var dictionary map[string][]int
var outputFileName string ="/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/animalcorpusoutput.txt"
var outputFile,_ = os.Create(outputFileName)
var totalDocs int
var totalWords int

func main() {

	//Scan input file and populate invertedIndex
	file := "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/animalcorpus.txt"
	InitializeInvertedIndex()
	totalWords = 0
	ScanFile(file)


	//First Prompt
	firstPrompt := []string{"Heidelberg" , "plutonium" , "Omarosa" , "octopus"}
	for i:= 0 ; i < len(firstPrompt) ; i++ {
		firstPrompt[i] = NormalizeText(firstPrompt[i])
	}
	//Print document frequency and postings list
	//postings list = (doc5, 3)
	printValues(firstPrompt , 1)


	//Second Prompt
	secondPrompt := []string{"Hopkins", "Harvard", "Stanford", "college"}
	for i:= 0 ; i< len(secondPrompt) ; i++ {
		secondPrompt[i] = NormalizeText(secondPrompt[i])
	}
	//Print document frequency only
	printValues(secondPrompt , 2)


	//Third Prompt
	thirdPrompt := []string{"Jeff" , "Bezos"}
	//Print out the docids that have both "Jeff" and "Bezos" in the text.
	for i:= 0 ; i < len(thirdPrompt) ; i++ {
		thirdPrompt[i] = NormalizeText(thirdPrompt[i])
	}
	findIntersection(thirdPrompt[0], thirdPrompt[1])

	//generate dictionary of word : offset , length
	InitializeDictionary()

	//create binary file output to write into
	outputFile.Close()

	//populate dictionary and append inverted index
	offset := 0
	i:= 0
	for key, value := range invertedIndex{
		offset = UpdateDictionary(key, value, offset)
		i++
	}

	//write dictionary into an output file
	WriteDictionary(dictionary)

	//Number of Documents
	fmt.Println( "Total Number of Documents : " + strconv.Itoa(totalDocs))

	//Size of dictioanry
	fmt.Println ("Total number of Unique Terms : " + strconv.Itoa(len(dictionary)))

	//Total number of Terms
	fmt.Println("Total number of Words in the Collection :" + strconv.Itoa(totalWords))

	fmt.Println("done")


}

func findIntersection(word_1 string, word_2 string){
	fmt.Println("Prompt #: 3")
	var sortedFirstKeys []int
	var sortedSecondKeys []int
	//get all of the docID that the first word is found in
	if firstkey, ok := invertedIndex[word_1]; ok {
		sortedFirstKeys = make([]int, 0, len(firstkey))
		for item := range firstkey {
			sortedFirstKeys = append(sortedFirstKeys, item)
		}
		sort.Ints(sortedFirstKeys)
	}
	//get all of the docID that the second word is found in
	if secondkey, ok := invertedIndex[word_2]; ok {
		sortedSecondKeys = make([]int, 0, len(secondkey))
		for item := range secondkey {
			sortedSecondKeys = append(sortedSecondKeys, item)
		}
		sort.Ints(sortedSecondKeys)
	}

	//find the intersection of the two arrays
	//determine which length is shorter
	anchor := sortedFirstKeys
	nonAnchor := sortedSecondKeys
	if len(sortedSecondKeys) < len(sortedFirstKeys) {
		anchor = sortedSecondKeys
		nonAnchor = sortedFirstKeys
	}
	var sharedDocs = make([]int, 0, len(anchor))
	i := 0
	j := 0
	for i < len(anchor){
		if anchor[i] == nonAnchor[j] {
			sharedDocs = append(sharedDocs, anchor[i])
			i++
			j++
		}else if anchor[i] > nonAnchor[j]{
			j++
		}else if anchor[i] < nonAnchor[j]{
			i++
		}


	}
	fmt.Println(sharedDocs)



}
func printValues (promptArray []string , promptNum int) {
	fmt.Println("Prompt #: " , promptNum)
	for i := 0 ; i < len(promptArray) ; i++ {
		fmt.Println (promptArray[i])
		if key, ok := invertedIndex[promptArray[i]]; ok {
			if promptNum == 1 {
				var keysToPrint []int
				i := 0
				for subkey, _ := range key {
					keysToPrint = Extend(keysToPrint, subkey)
					i++
				}
				sort.Ints(keysToPrint)
				for  i:=0 ; i < len(keysToPrint); i ++ {
					fmt.Println("docID : " , keysToPrint[i] , "doc freq : " , key[keysToPrint[i]])
				}

			} else {
				//determine the number of documents this word is in
				docFreq := len(key)
				fmt.Println ("document frequency : " , docFreq)
			}
		}else {
			fmt.Println("No document found for this term.")
		}
		fmt.Println()
	}
	fmt.Println()
}
//increase size of array and add element
func Extend(slice []int, element int) []int {
	n := len(slice)
	if n == cap(slice) {
		newSlice := make([]int, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
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
			}else if nextLine != "" { //process if current line is not blank. otherwise skip this line.
				//normalize text
				normalizedLine := NormalizeText(nextLine)
				//split up current line into an array of strings
				lineArray := s.Split(normalizedLine, " ")
				totalWords = totalWords + len(lineArray)
				//update inverted index
				UpdateIndex(lineArray,currentID)
			}
		}
	}
	totalDocs = currentID
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
	lineContent = s.Replace(lineContent,"@", " ",0 )
	lineContent = s.Replace(lineContent,"|", "",-1)
	lineContent = s.Replace(lineContent,"\\n", "",-1)
	lineContent = s.Replace(lineContent,"- ", " ",-2)
	lineContent = s.Replace(lineContent,"'s", "",-2)
	lineContent = s.Replace(lineContent,"'ve", "",-3)
	lineContent = s.Replace(lineContent,"'m", "",-2)
	lineContent = s.Replace(lineContent,"'re", "",-3)
	lineContent = s.Replace(lineContent,"'ed", "",-3)
	lineContent = s.Replace(lineContent,"'", "",-1)
	lineContent = s.Replace(lineContent,"  ", " ",-1)
	lineContent = s.Replace(lineContent,"  ", " ",-1)
	lineContent = s.Replace(lineContent,"!", "",-1)
	lineContent = s.Replace(lineContent,":", "",-1)
	return s.ToLower(lineContent)
}

//initialize invertedIndex
func InitializeInvertedIndex(){
	invertedIndex = make(map[string]map[int]int)
	return
}
//initialize dictionary
func InitializeDictionary(){
	dictionary = make(map[string][]int)
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
				}else { //create new item if docID does not exist in keyMap (postingsMap)
					key[docID] = 1
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
	newOffset := 0
	for key, value := range postingsMap {
		//get key values in binary
		keyBinary := make([]byte, 4)
		binary.LittleEndian.PutUint32(keyBinary, uint32(key))

		valueBinary := make([]byte, 4)
		binary.LittleEndian.PutUint32(valueBinary, uint32(value))

		WriteToFile(keyBinary)
		WriteToFile(valueBinary)
		newOffset = offset + len(keyBinary) + len(valueBinary)
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
	binary.Write(buf, binary.LittleEndian, input[3])
	binary.Write(buf, binary.LittleEndian, input[2])
	binary.Write(buf, binary.LittleEndian, input[1])
	binary.Write(buf, binary.LittleEndian, input[0])
	output := make([]byte, 4)
	output[0] = input[3]
	output[1] = input[2]
	output[2] = input[1]
	output[3] = input[0]
	if _, err = f.Write(output); err != nil {
		panic(err)
	}
}

func WriteDictionary (dictionary map[string][]int){
	var dictionaryFileName string ="/Users/mlo/GolandProjects/info-retrieval/Programming Assignments/dictionary.txt"
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