<template>
    <div class="modal" @click.self="$emit('close', Rename)">
    <div class="modal-content" @click.stop>
        <input v-model="Rename">
        <button v-if="ShowDeleteBtn" class="Delete" @click="DeleteFile">Delete File</button>
    </div>
    </div>
</template>

<script>
import axios from 'axios';
import helper from '../helper'
export default {
    props: ['items', 'FilePath'],
    data(){
        return {
            Rename : "",
            ShowDeleteBtn : ""
        }
    },
    created(){
        this.Rename = this.items
        let split = this.Rename.split(".")
        if(split && split.length > 1 && split[1]==='txt'){
            this.ShowDeleteBtn = true
        }
    },
    methods:{
        DeleteFile(){
            var self = this; 
            axios.get(`${helper.BaseUrl}deletefile?path=${this.FilePath+"/"+this.items}`)
            .then(function (response) {
                console.log(response);
                if(response.data.Status) {
                    alert("File Deleted Sucessfully")
                    self.$emit('close', self.items)
                }else{
                    alert(response.data.Message)
                }
            })
            .catch(function (error) {
                console.log(error);
            })
            .finally(function () {
            });
        }
    }
};
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

    .Delete {
        width: 10%;
        margin-top: 5px;
    }
</style>