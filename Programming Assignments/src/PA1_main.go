package main

import (
	"os"
	"bufio"
	s "strings"
	"github.com/google/btree"
	"fmt"
	"strconv"
)
//declare global tree creation
var list = btree.NewFreeList(32)
var tree = btree.NewWithFreeList(2,list)
var docID = 0
var docArray []string
var docToFind int
func main() {
	var file string
	var temp string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter full filepath to text document to read: ")
	file, _ = reader.ReadString('\n')
	file = s.Replace(file, "\n", "", -1)

	fmt.Print("Enter document ID to look up: ")
	temp, _ = reader.ReadString('\n')
	temp = s.Replace(temp, "\n", "", -1)
	docToFind, _ = strconv.Atoi(temp)

	// import in file into an array
	collectionFrequencyMap := make(map[string]int)
	documentFrequencyMap := make(map[string][]int)
	FileToLines(file, collectionFrequencyMap, documentFrequencyMap)

	//create reverse tree
	revCFMap := make(map[btree.Item][]string)
	createReverseMap(collectionFrequencyMap, revCFMap)

	//fill in binary tree
	updateTree(revCFMap)

	//get top 5 occurences
	var top100 [100]string
	var top100CF [100]int
	var top100DF [100][]int

	//top 500th
	var top500 string
	var top500CF int
	var top500DF []int

	k := 0
	for tree.Len() > 0 {
		output := tree.Max()
		strOutput := revCFMap[output]
		for j:= 0 ; j < len(strOutput) ; j++ {
			top100[k] = strOutput[j]
			top100CF[k] = collectionFrequencyMap[strOutput[j]]
			top100DF[k] = documentFrequencyMap[strOutput[j]]
			k++
			if k == 100 {
				break
			}
		}
		if k == 100 {
			break
		}else {
			tree.Delete(output)
		}
	}

	//Report the Number of Paragraphs processed
	fmt.Println ("Number of documents processed: " , docID)
	
	//Report the Number of of unique words processed
	fmt.Println("Number of of unique words processed: ", len(collectionFrequencyMap))

	//Report the Number of total words encountered
	var wordTotal int = 0
	for _, valueContent := range collectionFrequencyMap{
		wordTotal = wordTotal + valueContent
	}
	fmt.Println("Number of of total words processed: ", wordTotal)

	//Identify the 100 most frequent words (by total count, also known as collection frequency)
	// and report both the collection frequency and the document frequency for each.
	// Order them from most frequent to least.
	for i:=0 ; i < 100 ; i++ {
		fmt.Println("#" , i+1 , ": ", top100[i] , "\t collection frequency:" , top100CF[i] , " \t document frequency:" , len(top100DF[i]))
	}


	//500th most frequent word
	for tree.Len() > 0 {

		output := tree.Max()
		strOutput := revCFMap[output]
		for j:= 0 ; j < len(strOutput) ; j++ {
			if k == 500 {

				top500 = strOutput[j]
				top500CF = collectionFrequencyMap[strOutput[j]]
				top500DF = documentFrequencyMap[strOutput[j]]
				//fmt.Println(tree.Len(), "  " , k)
				k++
				break
			}else {
				k++
			}

		}
		if k == 501 {
			break
		}else {
			tree.Delete(output)
		}

	}
	fmt.Println("#500 : ", top500 , "\t collection frequency:" , top500CF , " \t document frequency:" , len(top500DF))

	//1000th most frequent word
	var top1000 string
	var top1000CF int
	var top1000DF []int
	for tree.Len() > 0 {
		output := tree.Max()
		strOutput := revCFMap[output]
		for j:= 0 ; j < len(strOutput) ; j++ {
			if k == 1000 {
				top1000 = strOutput[j]
				top1000CF = collectionFrequencyMap[strOutput[j]]
				top1000DF = documentFrequencyMap[strOutput[j]]
				k++
				break
			}else {
				k++
			}

		}
		if k == 1001 {
			break
		}else {
			tree.Delete(output)
		}
	}
	fmt.Println("#1000 : ", top1000 , "\t collection frequency:" , top1000CF , " \t document frequency:" , len(top1000DF))

	//1000th most frequent word
	var top5000 string
	var top5000CF int
	var top5000DF []int
	for tree.Len() > 0 {
		output := tree.Max()
		strOutput := revCFMap[output]
		for j:= 0 ; j < len(strOutput) ; j++ {
			if k == 5000 {
				top5000 = strOutput[j]
				top5000CF = collectionFrequencyMap[strOutput[j]]
				top5000DF = documentFrequencyMap[strOutput[j]]
				k++
				break
			}else {
				k++
			}

		}
		if k == 5001 {
			break
		}else {
			tree.Delete(output)
		}
	}
	fmt.Println("#5000 : ", top5000 , "\t collection frequency:" , top5000CF , " \t document frequency:" , len(top5000DF))


	//Report out the number of words that occur in exactly one document

	fmt.Println( "Number of words in document ID (" , docToFind , ") is " , WordsInDoc(docToFind, docArray))

	//Report out the percentage of dictionary items in a single document which is #uniquewords / #totalwords
	freqNumerator := UniqueWordsInDoc(docToFind, documentFrequencyMap)
	freqDenominator := WordsInDoc(docToFind, docArray)
	fmt.Println( "% of unique words in document ID (" , docToFind, ") is " , (float64(freqNumerator) / float64(freqDenominator)))



}

func FileToLines(filePath string, inputCFMap map[string]int , inputDFMap map[string][]int)  {


	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var textContent string
	for scanner.Scan() {
		nextLine := scanner.Text()
		if nextLine == "</P>" {
			docArray = ExtendString(docArray, textContent)
			textContent = ""
			docID++
		}else {
			if !s.Contains(nextLine, "<P"){
				nextLine = Normalize(nextLine)
				nextLine = RemovePunct(nextLine)
				textContent = textContent + nextLine
			}
		}
	}

	for i := 0 ; i < len(docArray) ; i++ {
		textArray := s.Split(docArray[i]," ")
		// create map
		// pass in map + document => update map -> new .go -> return updated map
		updateCFMap(textArray, inputCFMap)
		updateDFMap(i , textArray, inputDFMap)
	}}


func RemovePunct(lineContent string) (cleanedLine string){
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
	return
}

func Normalize(lineContent string)(normalizedLine string){
	normalizedLine = s.ToLower(lineContent)
	return
}

// check if value exists. if yes -> increment value by 1. if no -> add to map
func updateCFMap(text []string, wordMap map[string]int) {
	for j := 0; j < len(text) ; j++ {
		if key, ok := wordMap[text[j]]; ok {
			key++;
			wordMap[text[j]]=key;
		}else {
			wordMap[text[j]] = 1;
		}
	}
}
//check if current docID is listed in value. if yes -> do nothing. if no, add docID to list and increment counter by 1
func updateDFMap(docID int,text []string, wordMap map[string][]int){
	for j := 0; j < len(text) ; j++ {
		if key, ok := wordMap[text[j]]; ok { 	//check if word exists in the map
			docPresence := false 				// default set boolean to false
			for i:= 0; i<len(key) ; i++ {
				if docID == key[i] { 			//if docID exists in list already, set boolean to true
					docPresence = true;
				}
			}
			if docPresence == false { 			//if docID does not exist in list, append value to include docID
				key = Extend(key, docID)
				wordMap[text[j]] = key
			}
		}else {
			keyToAdd := []int{docID}		//if word doesn't exist in map, add word : docID
			wordMap[text[j]] = keyToAdd
		}
	}
}

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

func createReverseMap (inputMap (map[string]int), outputMap (map[btree.Item][]string)){
	for keyContent, valueContent := range inputMap{
		revKey := btree.Int(valueContent)
		revValue := keyContent
		if key, ok := outputMap[revKey]; ok { //if frequency number exists in map
			newValue := ExtendString(key, revValue)
			outputMap[revKey] = newValue
		}else{
			outputMap[revKey] = []string{revValue}
		}
	}
	return
}

func ExtendString(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func updateTree(inputMap (map[btree.Item][]string)) {
	for keyContent, _ := range inputMap{
		nodeToAdd := btree.Item(keyContent)
		tree.ReplaceOrInsert(nodeToAdd)
	}
	return
}

func WordsInDoc(docNum int ,  arrayDocs []string) int{
	tempArray := s.Split(arrayDocs[docNum] , " ")
	return len(tempArray)
}

func UniqueWordsInDoc (docNum int, DFmap map[string][]int) int{
	output := 0
	for _, valueContent := range DFmap{
		for j := 0 ; j < len(valueContent) ; j++ {
			if valueContent[j] == docNum {
				output++
			}
		}
	}
	return output
}