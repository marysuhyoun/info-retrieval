package main

import (
	"os"
	"bufio"
	"strconv"
	s "strings"
	"encoding/binary"
	"bytes"
	"math"
	"fmt"
	"sort"
	"time"
)

//query computation


var queryFile = "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14.topics.txt"
var dictionaryFile = "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14NoStemDictionary.txt"
var indexFile = "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/cds14NoStemInvertedIndex.txt"
var resultFile = "/Users/mlo/Documents/Personal/JHU/Information Retrieval/Module 4/Lo-a.txt"
var indexMap map[string]map[int]int
var termsMap map[string]int
var docMap map[int]int
var queryMap map[string]map[int]int
var stem string
var totalDocCount int
var scoreDoc map[int]map[string]float32 //map contains [docID] as int and [td-idf] score as float32
var scoreQuery map[int]map[string]float32 //map contains [queryID] as int and [td-idf] score as float32
var queryCosine [][]float64 //queryCosine is a 2D array to hold docID and cosine score
var queryDocOrder[][]int
func main(){
	fmt.Println("Query Start Time: ", time.Now())
	InitializeIndexMap()
	// create map from indexFile and dictionaryFile
	CreateIndex(indexFile, dictionaryFile)
	fmt.Println("Index File Read By Time: ", time.Now())
	totalDocCount = len(docMap)
	InitializeQueryMap()
	//read query text file and create a map of bag of words for each query.
	ReadQuery(queryFile)
	fmt.Println("Query File Read By Time: ", time.Now())
	//go through each document and query & assign scoring
	AssignTDIDF(indexMap, queryMap, totalDocCount)
	fmt.Println("td-idf Weights Assigned By Time: ", time.Now())
	//iterate through each query and generate cosine scores

	queryCosine := make([][]float64, len(scoreQuery))
	queryDocOrder = make([][]int, len(scoreQuery))
	for i := range queryCosine {
		queryCosine[i] = make([]float64, len(scoreDoc))
		queryDocOrder[i] = make([]int,len(scoreDoc))
	}
	fmt.Println("Cosine Scores Assigned By Time: ", time.Now())
	var tempScore float64

	for queryID, scoreMap := range scoreQuery{
		//iterate through each doc and get cosine similarity of every document
		i:= 0
		for documentID, documentMap := range scoreDoc{
			tempScore = CalculateCosine(queryID, scoreMap, documentMap)
			queryCosine[queryID-1][i] = tempScore
			queryDocOrder[queryID-1][i] = documentID
		}
	}
	fmt.Println("Cosine Similarity Calculated By  Time: ", time.Now())
	//Sort the queryCosine array by the cosine scoring high to low.
	sortedQueryCosine := make([][]float64, len(scoreQuery))
	for i := range sortedQueryCosine {
		sortedQueryCosine[i] = make([]float64, len(scoreDoc))
	}
	sortedQueryDoc := make([][]int, len(scoreQuery))
	for i := range sortedQueryDoc {
		sortedQueryDoc[i] = make([]int, len(scoreDoc))
	}
	SortCosineArray(queryCosine, queryDocOrder, sortedQueryCosine, sortedQueryDoc)
	fmt.Println("Documents Sorted By Cosine Similarity By Time: ", time.Now())
	//OUTPUTS
	//For first query, print out the query terms and their weights
	fmt.Println("Query #1 terms and weights")
	fmt.Println("Term \t TD/IDF Weight")
	firstQuery := scoreQuery[1]
	for queryTerm, queryScore := range firstQuery {
		fmt.Println(queryTerm, "\t", queryScore)
	}
	fmt.Println()

	//OUTPUT2 : Total size of unique words
	fmt.Println( "Total number of unique terms: " , len(indexMap))
	fmt.Println()
	//Output3 : Total number of documents in corpus
	fmt.Println( "Total number of documents in corpus: ", totalDocCount)



	//Create Output File with Sorted top 100 outputs
	PrintOutputFile(sortedQueryCosine, sortedQueryDoc)
	fmt.Println("Query End Time: ", time.Now())
	fmt.Println("done")
}

func PrintOutputFile(outputCosine [][]float64, outputDoc [][]int){
	resultFileName,_ := os.Create(resultFile)
	resultFileName.Close()
	f, err := os.OpenFile(resultFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	var max int
	max = 0
	for i:= 0; i < len(outputCosine); i++ {
		if len(outputCosine[i]) < 100 {
			max = len(outputCosine[i])
		}else{
			max = 100
		}
		for j:=0; j<max; j++{
			if _, err = f.WriteString(strconv.Itoa(i+1) + " Q0 " + strconv.Itoa(outputDoc[i][j]) + " " + strconv.Itoa(j+1) + " "  + strconv.FormatFloat(outputCosine[i][j],'f',4,64) + " Lo" + "\n"); err != nil {
			panic(err)
			}
		}

	}
	defer f.Close()

}
func SortCosineArray(inputArray[][]float64, documentOrder [][]int, outputScoreArray [][]float64 , outputDocArray[][]int){
	var docCosine map[int]float64 //lookup map for getting docID of the cosine scoring
	docCosine = make(map[int]float64)
	for i := 0 ; i < len(inputArray); i++{ //for every query
		for j := 0; j < len(inputArray[i]); j++{ //for every documentID
			docCosine[documentOrder[i][j]] = inputArray[i][j]
		}
		sort.Float64s(inputArray[i]) //sort the cosine scores.
		for j := 0; j < len(inputArray[i]); j ++ { //for every documentID
			outputScoreArray[i][j] = inputArray[i][len(inputArray[i]) - j - 1]
			for key,value := range docCosine{
				if value == outputScoreArray[i][j] {
					outputDocArray[i][j] = key
				}
			}
		}
	}
}
func CalculateCosine(queryID int, inputMap map[string]float32, docuMap map[string]float32)float64{
	var numerator float64
	numerator = 0.0
	for term,score := range inputMap{ //for every word in each query
		if _,ok := docuMap[term]; ok{
			numerator = numerator + (float64(docuMap[term]) * float64(score))
		}
	}
	var denominator float64
	var queryDenom float64
	var docDenom float64
	//denominator portion of query map
	for _,score := range inputMap {
		queryDenom = queryDenom + (float64(score) * float64(score))
	}
	queryDenom = math.Sqrt(float64(queryDenom))
	//denominator portion of document map
	for _,score := range docuMap {
		docDenom = docDenom + (float64(score) * float64(score))
	}
	docDenom = math.Sqrt(float64(docDenom))
	denominator = queryDenom * docDenom
	var valueToReturn float64
	valueToReturn= 0
	if denominator == 0.0 {
		valueToReturn = 0
	}else{
		valueToReturn = numerator/denominator
	}
	return valueToReturn
}

func AssignTDIDF(inputMap map[string]map[int]int, queryMap map[string]map[int]int, totalDocNum int){
	//TDIDF = term frequency in document/query * log2(total # doc/queries / document/query frequency)
	//initialize score maps
	scoreDoc = make(map[int]map[string]float32)
	scoreQuery = make(map[int]map[string]float32)
	//iterate through indexMap and calculate TD and IDF. Populate scoreDoc with TD*IDF
	for uniqueTerm, postingMap := range inputMap {
		//length of each posting map is the total number of documents that contain this term (denominator of IDF)
		dF := len(postingMap)
		//loop through each document associated with this term and place the term frequency inside the map
		for docNum, termFreq := range postingMap{
			if key, ok := scoreDoc[docNum]; ok {
				key[uniqueTerm] = float32(float64(termFreq) * math.Log2(float64(totalDocNum)/float64(dF)))
			}else{
				wordMap := make(map[string]float32)
				wordMap[uniqueTerm] =float32(float64(termFreq) * math.Log2(float64(totalDocNum)/float64(dF)))
				scoreDoc[docNum] = wordMap
			}
		}
		if key, ok := queryMap[uniqueTerm]; ok{
			for queryNum, wordFreq := range key{
				if key, ok := scoreQuery[queryNum]; ok {
					key[uniqueTerm] = float32(float64(wordFreq) * math.Log2(float64(totalDocNum)/float64(dF)))
				}else{
					wordMap := make(map[string]float32)
					wordMap[uniqueTerm] =float32(float64(wordFreq) * math.Log2(float64(totalDocNum)/float64(dF)))
					scoreQuery[queryNum] = wordMap
				}
			}
		}
	}
}

func InitializeIndexMap(){
	indexMap = make(map[string]map[int]int)
	docMap = make(map[int]int)
	return
}

func InitializeQueryMap(){
	queryMap = make(map[string]map[int]int)
	return
}
func CreateIndex(indexFile string, dictionaryFile string) {
	f, err := os.Open(dictionaryFile)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var inputArray []string
	var length int
	var offset int

	scanner.Scan() //skip first line (header)
	for scanner.Scan(){
		nextLine := scanner.Text()
		inputArray = s.Split(nextLine, " , ") //[term, length, offset]
		length, _ = strconv.Atoi(inputArray[1])
		offset, _ = strconv.Atoi(inputArray[2])
		AddToIndex(inputArray[0], length, offset, indexFile)
		if _, ok := termsMap[inputArray[0]]; ok {
			termsMap[inputArray[0]] = termsMap[inputArray[0]] + length
		}else{
			termsMap = make(map[string]int)
			termsMap[inputArray[0]] = length
		}
	}
	f.Close()
}

//add term to indexMap
func AddToIndex(term string, length int, offset int, indexFile string){
	g, err := os.Open(indexFile)
	if err != nil {
		return
	}
	defer g.Close()
	reader := bufio.NewReader(g)
	_,err = reader.Discard(offset)
	p := make([]byte, length*8)
	reader.Read(p)
	i := 0
	var docID int32
	var docFreq int32
	docIDByteArray:= make([]byte,4)
	docFreqByteArray:= make([]byte,4)
	postingsMap := make(map[int]int)
	indexMap[term] = make(map[int]int)
	for i < len(p) {
		for j := 0 ; j < 4 ; j++ {
			docIDByteArray[j] = p[j+i]
			docFreqByteArray[j] = p[j+i+4]
		}
		buf := bytes.NewReader(docIDByteArray)
		binary.Read(buf, binary.LittleEndian, &docID)
		buf2 := bytes.NewReader(docFreqByteArray)
		binary.Read(buf2, binary.LittleEndian, &docFreq)
		i = i + 8
		postingsMap[int(docID)] = int(docFreq)

		if _,ok := docMap[int(docID)];ok {
		}else{
			docMap[int(docID)] = 1
		}
	}
	indexMap[term] = postingsMap
	if math.Mod(float64(len(indexMap)),500.0) == 0 {
		fmt.Println(len(indexMap))
	}
	g.Close()
}

//read query file and create map of bag of words
func ReadQuery(queryFile string){
	f, err := os.Open(queryFile)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var currentID int
	for scanner.Scan(){
		nextLine := scanner.Text()
		// </P> marks the end of the document.
		if !s.Contains(nextLine, "</Q>"){
			//determine docID
			if s.HasPrefix(nextLine, "<Q ID="){
				//remove prefix and suffix of docID
				nextLine = s.TrimPrefix(nextLine, "<Q ID=")
				nextLine = s.TrimSuffix(nextLine, ">")
				currentID,_ = strconv.Atoi(nextLine)
			}else if nextLine != "" { //process if current line is not blank. otherwise skip this line.
				//normalize text (non-stemming)
				normalizedLine := NormalizeLine(nextLine)
				//remove stop words from the list. https://gist.github.com/sebleier/554280
				normalizedLine = DeleteStopWords(normalizedLine)
				//split up current line into an array of strings
				lineArray := s.Split(normalizedLine, " ")
				//stem if required
				if stem == "yes" {
					StemLine(lineArray)
				}
				//update inverted index
				UpdateQueryMap(lineArray,currentID)
			}
		}
	}
}

//update Query Map with query line
func UpdateQueryMap(lineContents []string, queryID int){
	//make a termMap for every unique word in the queryMap
	for i := 0 ; i < len(lineContents) ; i++ {
		//make sure the string in question is not composed of just " " spaces.
		lookUpItem := s.TrimSpace(lineContents[i])
		if lookUpItem != "" {
			//see if term already exists in queryMap
			if key, ok := queryMap[lineContents[i]]; ok{
				//see if term already exists in keyMap (termMap)
				if _, ok := key[queryID]; ok {
					key[queryID] = key[queryID] + 1
				}else { //create new item if term does not exist in keyMap (termMap)
					key[queryID] = 1
				}
			} else { //create a new tuple item if one does not exist
				termMap := make(map[int]int)
				termMap[queryID] = 1
				queryMap[lineContents[i]] = make(map[int]int)
				queryMap[lineContents[i]] = termMap
			}
		}

	}
}
//given set of text, normalize it
func NormalizeLine(lineContent string) string{
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
func DeleteStopWords(lineContent string) string {
	stopwords := [127]string{"i","me","my","myself","we","our","ours","ourselves","you","your","yours","yourself","yourselves","he","him","his","himself","she","her","hers","herself","it","its","itself","they","them","their","theirs","themselves","what","which","who","whom","this","that","these","those","am","is","are","was","were","be","been","being","have","has","had","having","do","does","did","doing","a","an","the","and","but","if","or","because","as","until","while","of","at","by","for","with","about","against","between","into","through","during","before","after","above","below","to","from","up","down","in","out","on","off","over","under","again","further","then","once","here","there","when","where","why","how","all","any","both","each","few","more","most","other","some","such","no","nor","not","only","own","same","so","than","too","very","s","t","can","will","just","don","should","now"}
	for i := 0 ; i < len(stopwords) ; i++{
		if s.Contains(lineContent, stopwords[i]){
			s.Replace(lineContent, stopwords[i],"", -1)
		}
	}
	return lineContent
}
//given set of text, stem it to the first 5 characters
func StemLine(lineArray []string) []string{
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