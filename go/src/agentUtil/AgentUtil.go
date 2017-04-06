

// version No 1 dated :- 03-Apr-2017
package agentUtil
// version No 1 dated :- 03-Apr-2017
import (
    
    "os/exec"
  _ "fmt" // for unused variable issue
    "fileUtil"
     "fmt"
    "io/ioutil"
     "encoding/json"
     "net/http"
     "stringUtil"
      "strconv"
      "strings"
      //"reflect"
   
)

//const url = "https://8zq40bjoqc.execute-api.us-west-2.amazonaws.com/dev/hello"
//const url = "https://xhas8lgghd.execute-api.us-west-2.amazonaws.com/dev/getinstructionsforagent?serverIp=192.168.1.1"
const url = "https://xhas8lgghd.execute-api.us-west-2.amazonaws.com/dev/getinstructionsforagent?serverIp=172.31.15.1"
var maxSize int = 100
var array = make([]string, maxSize) 
var cntr int = 0
const Delimiter = ":"


func ExecComand(cmd, fromFile string) string {
    
    cmdStatus,err := exec.Command("bash","-c",cmd).Output()
    execStatus := "success"
    if err != nil {
        errorMsg := "Cmd executed = : "+cmd +" : execStatus = : fail. fromFile. = :"+fromFile
        fileUtil.WriteIntoLogFile(errorMsg)
        execStatus = "fail"
        fmt.Println("34. AgentUtil.ExecComand()  errorMsg = : ", errorMsg)

    }

    if (len(string(cmdStatus)) > 0){
        execStatus =  string(cmdStatus)  
    }
    return execStatus
}

func hitApi() string{
  
    array = make([]string, maxSize)
    cntr = 0

    fileUtil.WriteIntoLogFile("UserHandler.getNextWork(). Going to hit url = : "+url)
    //fmt.Println("UserHandler.getNextWork(). Going to hit url = : "+url)
    res, err := http.Get(url)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN 119. Msg = : "+err.Error())
        return "1"
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
      fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN. 124. Msg = : "+err.Error())
      return "1"
    }
     
    var data []interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        fileUtil.WriteIntoLogFile("Error at UserHandler.getNextWork(). LN 130. Msg = : "+err.Error())
        return "1"
    }
    //fmt.Println(data)
    //parseArray(data)
    parseArray(data)
    TrimArrayToActualSize()
    return "0"
    
   
}//hitApi



func GetNextWork() ([]string){
  for i := 0; i < 10; i++ {
    //var infraGuardResponse []string
    respStatus := hitApi();
    if(respStatus == "0"){
     
      isValidData := ValidateArray()
     
      msg := "\n************ Is server responses valid data =  "+strconv.FormatBool(isValidData) +" ***********************\n"
      fileUtil.WriteIntoLogFile(msg)
      if(isValidData){
        
        fmt.Println(msg)
         printArray()

         if(len(array) > 0 ){
             
             var tmp = array
             array = nil
             return tmp 
         }else{
          return nil;
         }
        
      }/*else{
        printArray()
      }*/
    }
    
  }
  return nil
}

func parseMap(aMap map[string]interface{}) {
    for key, val := range aMap {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
           // initializeArray(key, "")
            parseMap(val.(map[string]interface{}))

        case []interface{}:
            //initializeArray(key, "")
            parseArray(val.([]interface{}))
        default:
            if val, ok := concreteVal.(string); ok {
               initializeArray(key, val)
             
          
            }

           
        }
    }
}

func parseArray(anArray []interface{}) {
    for _, val := range anArray {
        switch concreteVal := val.(type) {
        case map[string]interface{}:
            parseMap(val.(map[string]interface{}))
        case []interface{}:
            parseArray(val.([]interface{}))
        default:
            if val, ok := concreteVal.(string); ok {
               initializeArray(val, "")
            }

        }
    }
}


func initializeArray(key, val string){
  key = stringUtil.RemoveSpace(key)
  val = stringUtil.RemoveSpace(val)
  if(key == "reqData" || val == "reqData" || len(key) == 0 || len(val) == 0){
    return
  }
  if(key == "requiredData"){
     var values []string
     values = stringUtil.SplitData(val, ",")
      for i := 0; i < len(values); i++ {
           values[i] = strings.Replace(values[i], "{", "", -1)
           values[i] = strings.Replace(values[i], "}", "", -1)
           values[i] = strings.Replace(values[i], "\"", "", -1)
           array[cntr] = values[i]
           cntr = cntr + 1
      }
      return
  }
  array[cntr] = key + Delimiter + val
  cntr = cntr + 1
  return
}

func printArray(){
  fmt.Println("Agent has below following works")
  fileUtil.WriteIntoLogFile("Agent has below following works")
   for i := 0; i < len(array); i++ {
    val := array[i]
    fmt.Println(val)
    fileUtil.WriteIntoLogFile(val)
  }
}



func ValidateArray() bool{
  
  var values []string
 //printArray()

   for i := 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
     
  
   if(values[0] == "activityName"){
      isValidData := checkDataSequence(values[1], i)
      if(isValidData == false){
        return false;
      }
   }

   
  }

  isValidData := checkActivityNameSequence()
  return isValidData
}



func checkDataSequence(activityName string, cnt int) bool{

    var values, sequnce []string
  
    if(activityName == "addUser"){
      sequnce = []string{"activityName","publicKey", "userName", "shell", "id"}
    }

    if(activityName == "deleteUser"){
      sequnce = []string{"activityName","userName","id"}
    }
    if(activityName == "changePrivilege"){
      sequnce = []string{"activityName","userName", "privilege", "id"}
    }
    
    if((cnt + len(sequnce)) > len(array)){
      fmt.Println("Returning false from 209 ")
      return false
    }
   

     for i := 0; i < len(sequnce); i++{
        values = stringUtil.SplitData(array[cnt], Delimiter)
        cnt++
  
        if(strings.ToLower(values[0]) != strings.ToLower(sequnce[i])){
           return false
        }
      
      }
    return true 
  }//checkDataSequence


func checkActivityNameSequence() bool{
  var values []string
  var i, cnt int
  cnt = 0
  for i = 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
   if(values[0] == "activityName"){
     cnt++;
   }

  }
  var tmpArray = make([]string, cnt)
  cnt = 0
  
  for i = 0; i < len(array); i++ {
   values = stringUtil.SplitData(array[i], Delimiter)
   if(values[0] == "activityName"){
     tmpArray[cnt] = strconv.Itoa(i) + ":"+values[1]
     cnt++;
   }

  }


  for j := 0; j < len(tmpArray)-1; j++ {
    
     values = stringUtil.SplitData(tmpArray[j], Delimiter)
     val1, _ := strconv.ParseInt(values[0], 10, 0)
     activityName1 := values[1]
     
     var jInt64 int64
     jInt64 = int64(j)
     if(jInt64 == 0 && jInt64 != val1){
      return false
     }

     values = stringUtil.SplitData(tmpArray[j+1], Delimiter)
     val2, _ := strconv.ParseInt(values[0], 10, 0)
    

     if(activityName1 == "addUser"){
       if((val1 + 5 )!= val2){
        return false
       }
     }

     if(activityName1 == "deleteUser"){
       if((val1 + 3 )!= val2){
        return false
       }
     }
     if(activityName1 == "changePrivilege"){
       if((val1 + 4 )!= val2){
        return false
       }
     }
  }//I
 return true



}

func TrimArrayToActualSize(){
  var cnt int = 0 
  for i := 0; i < len(array); i++ {
  val := array[i]
    if(len(val) > 0){
      cnt++
    }
  }
  var tmp = array[0:cnt]
  array = tmp
  tmp = nil
 
  
}

func HitAnyUrl(serverUrl string) string{
  res, err := http.Get(serverUrl)
  if err != nil {
      fileUtil.WriteIntoLogFile("Error at AgentUtil.hitAnyUrl(). LN 297. Msg = : "+err.Error())
      return "1"
  }
  _, error := ioutil.ReadAll(res.Body)
  if error != nil {
    fileUtil.WriteIntoLogFile("Error at AgentUtil.hitAnyUrl(). LN 305.. Msg = : "+error.Error())
    return "1"
  }

  fileUtil.WriteIntoLogFile("Successfully sent execution status to this url = : "+serverUrl)
  fmt.Println("Successfully sent execution status to this url = : ",serverUrl) //success
  return "0"
}

