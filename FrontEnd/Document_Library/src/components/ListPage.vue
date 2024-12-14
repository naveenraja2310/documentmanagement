
<template>
  <!--Sub Navbar to Display Path and Required Buttons-->
  <div class="subNavBar">
    <button  type="button" class="btn btn-danger mx-2"  @click="BackButton()" :disabled="RoutePath == ''" style="width: 10%;">Back</button>
    <button type="button" class="btn btn-primary mx-2"  @click="ShowTextFile()" style="width: 10%; margin-left: 10px;">Add Text File</button>
    <input :value="ValueinSearchBar" disabled style="width: 80%; margin-left: 10px;">
  </div>

  <!--Displaying Folders and Files-->
  <div>
    <ul class="file-list">
      <li v-for="ff in FolderandFiles" :key="ff.name" class="file-item">
        <div @contextmenu.prevent="openContextMenu(ff)">
          <span v-if="ff.isDir">
            <i class="fas fa-folder fa-3x folder-icon" @dblclick="OnClickFolder(ff.name)"></i>
          </span>
          <span v-else>
            <i class="fas fa-file fa-3x" @click="DownloadFile(ff.name)"></i>
          </span>
          <div class="file-name">{{ ff.name }}</div>
        </div>
      </li>
    </ul>
  </div>

  <!--Space for Components-->
  <Modal v-if="showContextMenu" 
    @close="closeContextMenu" 
    :items="contextMenuItems"
    :FilePath="RootPath + RoutePath"
    
  ></Modal>

  <TextFile v-if="showTextFile" 
    @closeTextFile="closeTextFile"
    :FilePath="RootPath + RoutePath"
    :textContent="textContent"
    :NameofFile="NameofFile"
  ></TextFile>
</template>

<script>

import axios from 'axios';
import Modal from './Modal.vue'
import TextFile from './TextFile.vue';
import helper from '../helper'
import "@fortawesome/fontawesome-free/css/all.css";


export default{
  data() {
    return {
      FolderandFiles : [],
      RootPath : "../",
      RoutePath: "",
      showContextMenu: false,
      contextMenuItems: "",
      showTextFile:false,
      textContent:"",
      NameofFile :"",
    }
  },
  computed:{
    ValueinSearchBar() {
      return this.RoutePath
    }
  },
  created(){
    this.getlist()
  },
  components:{
    Modal,
    TextFile
  },
  methods:{
    getlist(){
      var self = this; 
      axios.get(`${helper.BaseUrl}list?path=${this.RootPath}`)
      .then(function (response) {
        console.log(response);
        if(response.data.Status) {
          self.FolderandFiles = response.data.Data
        }else{
          alert(response.data.Message)
        }
      })
      .catch(function (error) {
        console.log(error);
      })
      .finally(function () {
      });
    },

    OnClickFolder(name, myPath){
      var self = this; 
      let URL = ""
      if(myPath || myPath === '') {
        URL = `${helper.BaseUrl}list?path=${this.RootPath+myPath}`
      }else{
        self.RoutePath += "/"+name
        URL = `${helper.BaseUrl}list?path=${this.RootPath+self.RoutePath}`
      }
      axios.get(URL)
      .then(function (response) {
        console.log(response);
        if(response.data.Status) {
          self.FolderandFiles = response.data.Data
        }else{
          alert(response.data.Message)
        }
      })
      .catch(function (error) {
        console.log(error);
      })
      .finally(function () {
      });
    },

    DownloadFile(fileName) {
      const self = this;
      let splitFileName = fileName.split(".")
      if(splitFileName && splitFileName.length>1 && splitFileName[1]==='txt'){
        this.readTextFile(fileName)
      }else{
        const URL = `${helper.BaseUrl}download?path=${self.RootPath + "/" + self.RoutePath +"/"+ fileName}`;
        axios.get(URL, { responseType: 'blob' })
          .then(response => {
            const url = window.URL.createObjectURL(new Blob([response.data]));
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', fileName);
            document.body.appendChild(link);
            link.click();
            document.body.removeChild(link);
          })
          .catch(error => {
            console.error('Error downloading file:', error);
            alert('Error downloading file. Please try again.');
          });
      }      
    },

    readTextFile(fileName) {
      console.log("FileNameForModal", fileName)
      var self = this;
      axios.get(`${helper.BaseUrl}readfile?path=${self.RootPath + "/" + self.RoutePath +"/"+ fileName}`)
        .then(response => {
          if(response.data.Status){
      console.log("FileNameForModal", fileName)

            self.textContent = response.data.Data;
            self.NameofFile = fileName
            self.showTextFile = true
          }else{
            alert(response.data.Message) 
          }
        })
        .catch(error => {
          console.error('Error reading file:', error);
          alert('Failed to read file');
        });
    },

    BackButton(){
      console.log(this.RoutePath)
      const pathSegments = this.RoutePath.split('/');
      pathSegments.pop();
      this.RoutePath = pathSegments.join('/');
      this.OnClickFolder(undefined, this.RoutePath);
    },

    openContextMenu(file) {
      this.contextMenuItems = file.name
      this.showContextMenu = true;
    },

    closeContextMenu(newName) {
      console.log("new name", newName)
      if(this.contextMenuItems != newName){
        var self = this; 
        let URL = `${helper.BaseUrl}rename?oldPath=${self.RootPath + "/" + self.RoutePath +"/"+this.contextMenuItems}&newPath=${self.RootPath + "/" + self.RoutePath +"/"+newName}`
        axios.get(URL)
        .then(function (response) {
          console.log(response);
          if(response.data.Status) {
            self.showContextMenu = false;
            self.contextMenuItems = "";
            self.OnClickFolder(undefined, self.RoutePath)
          }else{
            alert(response.data.Message)
          }
        })
        .catch(function (error) {
          console.log(error);
        })
        .finally(function () {
        });
      }else{
        this.showContextMenu = false;
        this.contextMenuItems = "";
        this.OnClickFolder(undefined, this.RoutePath)
      }
    },

    ShowTextFile(){
      this.showTextFile = true
    },

    closeTextFile(){
      this.showTextFile = false
      this.textContent=""
      this.NameofFile = ""
      this.OnClickFolder(undefined, this.RoutePath)
    }
    
  }
}

</script>


<style scoped>
ul li h3 {
    list-style-type: none;
    margin: 0;
    padding: 0;
    overflow: hidden;
    background-color: #333;
    display: flex; /* Use flexbox for layout */
    justify-content: space-between; /* Distribute items evenly across the width */
    }
.file-list {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  list-style-type: none;
  padding: 0;
  margin-top: 20px;
}

.file-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 120px;
  margin: 20px;
  text-align: center;
}

.file-name {
  margin-top: 10px;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.subNavBar {
  display: flex;
  padding: 2.5px;
  background-color: beige;
}
.folder-icon {
  color: #ffd04b; /* Blue color for folder icon */
}
</style>
