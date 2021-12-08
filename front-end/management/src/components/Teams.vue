<template>
  <div class="table">
    <meta http-equiv="Content-Security-Policy" content="default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com">
    <h1>LIST OF TEAMS</h1>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">

    <button class="buttonAdd" onclick="document.getElementById('addModal').style.display='block'"><i class="fa fa-plus"></i> ADD NEW TEAM</button>

    <div id="addModal" class="w3-modal">
      <div class="w3-modal-content">
      <div class="w3-container">
        <form class="form">
          <h2>ADD NEW TEAM</h2>
          <br>
          <label>Name</label>&nbsp;
          <input type="text" v-model="team.name" placeholder="Enter name..." required="required">
          <br><br>
          <label>Description</label>&nbsp;
          <input type="text" v-model="team.description" placeholder="Enter description..." required="required">
          
          </form>
        <button class="button" @click="this.postItem();" onclick="document.getElementById('addModal').style.display='none';">Register</button>
        <button class="button" onclick="document.getElementById('addModal').style.display='none'">Cancel</button>
      </div>
      </div>
    </div>

    <div id="editModal" class="w3-modal">
    <div class="w3-modal-content">
      <div class="w3-container">
        <form>
          <h2>EDIT A TEAM</h2>
          <br>
          <label>[NAME] Old: [ {{this.team.name}} ] &nbsp;&nbsp;&nbsp; New: </label>&nbsp;
          <input type="text" v-model="newTeam.name" placeholder="New name..." required>
          <br><br>
          <label>[Description] Old: [ {{this.team.description}} ] &nbsp;&nbsp;&nbsp; New: </label>&nbsp;
          <input type="text" v-model="newTeam.description" placeholder="New description..." required="required">
          <br><br>

          <button class="button" @click="this.editItem(this.uid)" onclick="document.getElementById('editModal').style.display='none'">Edit</button>
          <button class="button" onclick="document.getElementById('editModal').style.display='none'">Cancel</button>
          </form>
        
      </div>
    </div>


    </div>

    <div id="delModal" class="w3-modal">
    <div class="w3-modal-content">
      <div class="w3-container">
          <h2>DELETE A TEAM</h2>
          <br>
          <label>Name:  {{this.team.name}} </label>
          <br><br>
          <label>Description: {{this.team.description}}  </label>
          <br><br>
          <button class="button" @click="this.deleteItem(this.uid)" onclick="document.getElementById('delModal').style.display='none'">Delete</button>
          <button class="button" onclick="document.getElementById('delModal').style.display='none'">Cancel</button>      
      </div>
    </div>
    </div>

    <table class="table">
      <thead>
        <tr>
          <th class="id">ID</th>
          <th class="name">NAME</th>
          <th class="description">DESCRIPTION</th>
          <th class="option">OPTIONS</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index)  in list" :key="item.uid">
          <th class="id">{{++index}}</th>
          <th class="name">{{item.name}}</th>
          <th class="description">{{item.description}}</th>
          <th class="option">
            <button class="button" @click="this.$router.push({name:'Team', params:{uid: item.uid}})"><i class="fa fa-search"></i> Detail</button>
            <button class="button" @click="this.team=item; this.uid=item.uid" onclick="document.getElementById('editModal').style.display='block'"><i class="fa fa-edit"></i> Edit</button>
            <button class="button" @click="this.team=item; this.uid=item.uid" onclick="document.getElementById('delModal').style.display='block'"><i class="fa fa-trash"></i> Delete</button>
          </th>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>


import axios from "axios"
export default {
  name: 'B',
  data () {
    return {
      list: null,
      team:{
        name: null,
        description: null,
      },
      newTeam:{
        name: null,
        description: null,
      },
      uid: null,
    }
  },
  methods: {
    getItem(){
      axios.get("http://localhost:8081/api/v1/team/")
      .then((response) => {
        this.list = Array.isArray(response.data.data)? response.data.data: [];
        console.log(this.list);
      }).catch(error => console.error(error));
      

    },

    postItem(){
      axios.post("http://localhost:8081/api/v1/team/", this.team).then((res) => {
        this.list.push(res.data.data);
      }).catch(error => console.error(error));
      // location.reload();
    },

    editItem(uid){
      axios.patch(`http://localhost:8081/api/v1/team/${uid}`, this.newTeam).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      alert("Updating...")
      var delayInMilliseconds = 1000; 
      setTimeout(function() {
        // location.reload();
      }, delayInMilliseconds);
    },

    deleteItem(uid){
      axios.delete("http://localhost:8081/api/v1/team/"+uid).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      alert("Deleting...")
      var delayInMilliseconds = 1000; 
      setTimeout(function() {
        location.reload();
      }, delayInMilliseconds);
    }

  },
  created() {
    this.getItem();
  }


}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">

table, td, th {
  border: 1px solid green;
}

table {
  border: 1px solid green;
  width: 80%;
  margin: auto;
  margin-bottom: 100px;
}

thead {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: center;
  background-color: #4CAF50;
  color: white;
}

.id{
  width: 5%;
}

.name{
  width: 15%;
}

.description{
  width: 30%;
}

.option{
  width: 30%;
}

.form{
  width: 378px;
  height: 256px;
  align-content: center;
  margin: auto;
}


tr {
  height: 40px;
  vertical-align: center;
}

.button {
  background-color: #4CAF50; /* Green */
  border: none;
  color: white;
  padding: 16px 32px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
  font-size: 16px;
  margin: 4px 2px;
  transition-duration: 0.4s;
  cursor: pointer;
}
.button:hover {
  background-color: #4CAF0A;
  color: black;
}

.buttonAdd {
  background-color: #4CAF50; /* Green */
  border: black;
  color: white;
  padding: 16px 32px;
  text-align: center;
  text-decoration: none;
  font-size: 16px;
  margin: 4px 2px;
  position: relative;
  right: -35%;
  transition-duration: 0.4s;
  cursor: pointer;
}
.buttonAdd:hover {
  background-color: #4CAF0A;
  color: black;
}

thead {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: center;
  background-color: #4CAF50;
  color: white;
}

</style>