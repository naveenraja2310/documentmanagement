<template>
    <div class="modal">
        <div class="modal-content">
            <form>
                <div class="form-group" v-if="!ModeUpdate">
                    <label for="exampleInputEmail1">File Name {{ FileName }}</label>
                    <input type="text" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp" v-model="FileName" placeholder="File Name">
                </div>
                <div class="form-group">
                    <label for="exampleInputPassword1">File Content</label>
                    <textarea type="text" class="form-control" id="exampleInputPassword1" placeholder="File Content" v-model="TextContent" rows="10"></textarea>
                </div>
                <div class="row justify-content-center my-3">
                    <div class="col-md-auto d-flex justify-content-between">
                        <button type="button" class="btn btn-danger mx-2"  @click="$emit('closeTextFile')">Cancel</button>
                        <button type="button" class="btn btn-primary mx-2" @click="createTextFile" v-if="!ModeUpdate">Save File</button>
                        <button type="button" class="btn btn-primary mx-2" @click="updateTextFile" v-else>Update File</button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import helper from '../helper'
export default {
    data(){
        return {
            FileName:"",
            TextContent : '',
            ModeUpdate : false
        }
    },
    props: ['FilePath', 'textContent', 'NameofFile'],
    created(){
        console.log("this.FileName1", this.NameofFile)
        if(this.NameofFile){
            this.ModeUpdate = true
        }
        this.FileName = this.NameofFile
        this.TextContent = this.textContent
    },
    methods:{
        createTextFile() {
            var self = this
            var fname = self.RemoveSpaceandSpecialChar(this.FileName)
            axios.post(`${helper.BaseUrl}createtextfile?path=${this.FilePath}&fileName=${fname+'.txt'}`, this.TextContent)
            .then(response => {
            console.log(response.data);
            alert('Text file created successfully');
                self.FileName = '';
                self.TextContent = '';
                self.$emit('closeTextFile')
            })
            .catch(error => {
                console.error('Error creating text file:', error);
                alert('Failed to create text file');
            });
        },
        updateTextFile() {
            var self = this
            debugger
            var path = `${self.FilePath + "/"+ self.FileName}`
            axios.post(`${helper.BaseUrl}updatefile`, new URLSearchParams({
                path: path,
                content: self.TextContent
            }))
            .then(response => {
                console.log(response.data);
                if(response.data.Status){
                    alert('File updated successfully');
                    self.FileName = '';
                    self.TextContent = '';
                    self.$emit('closeTextFile')
                }else{
                    alert(response.data.Message);
                }
            })
            .catch(error => {
            console.error('Error updating file:', error);
            alert('Failed to update file');
            });
        },

        RemoveSpaceandSpecialChar(newValue) {
            newValue = newValue.replace(/[^\w\s]/gi, '')
            newValue = newValue.replace(/\s+/g, '')
            return newValue.toLowerCase()
        },
    }

}
</script>

<style scoped>
    .modal {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    }

    .modal-content {
    background-color: #fff;
    padding: 20px;
    border-radius: 5px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.3);
    max-width: 80%;
    max-height: 80%;
    overflow: auto;
    }

    .modal-item {
    cursor: pointer;
    padding: 10px;
    transition: background-color 0.3s ease;
    }

    .modal-item:hover {
    background-color: #f0f0f0;
    }
</style>