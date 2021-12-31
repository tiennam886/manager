<template>
    <body>
      <div class="all" v-if="this.team">
        <meta http-equiv="Content-Security-Policy" content="default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com">
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
        <div class="card" >
          <div><img src="../assets/logo.png" alt="Avatar" style="width:100%"></div>
          <div class="container">
          <h4><b>{{this.team.name}}</b></h4>  
          <p>{{this.team.description}}</p> 
          </div>
          
      </div>

      <div>
        <div class="add">
          <button class="button" onclick="document.getElementById('addModal').style.display='block'"><i class="fa fa-plus"></i> ADD NEW EMPLOYEE</button>
          <button class="button" @click="this.update(this.team);" onclick="document.getElementById('editModal').style.display='block'"><i class="fa fa-edit"></i> Edit</button>
          <button class="button" onclick="document.getElementById('delModal').style.display='block'"><i class="fa fa-trash"></i> Delete</button>
              
              <div id="addModal" class="w3-modal" >
                <div class="w3-modal-content">
                <div class="w3-container" v-if="employees.length!=0">
                <table class="table">
                  <thead>
                    <tr>
                      <th class="id">ID</th>
                      <th class="name">NAME</th>
                      <th class="gender">GENDER</th>
                      <th class="dob">BIRTH</th>                      
                      <th class="option">OPTIONS</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item, index) in employees" :key="item.uid">
                      <th class="id">{{++index}}</th>
                      <th class="name">{{item.name}}</th>
                      <th class="gender">{{item.gender}}</th>
                      <th class="dob">{{item.dob}}</th>                      
                      <th class="option">
                        <button class="button" @click="this.chosenEmployee=item;" onclick="document.getElementById('confirmAddModal').style.display='block'">
                          <i class="fa fa-plus"></i> Add</button>
                        <div id="confirmAddModal" class="w3-modal">
                          <div class="w3-modal-content">
                            <h2> <b>ADD NEW EMPLOYEE</b> </h2>
                
                            <h5>Do you want to add <b>{{chosenEmployee.name}} to this team?</b></h5>
                            <br>
                            <button class="button" @click="this.addToTeam(team.uid, chosenEmployee.uid)" 
                            onclick="document.getElementById('confirmAddModal').style.display='none'; document.getElementById('addModal').style.display='none'">Confirm</button>
                            <button class="button" onclick="document.getElementById('confirmAddModal').style.display='none'">Cancel</button>                            
                          </div>
                        </div>
                      </th>
                    </tr>
                  </tbody>
                </table>                 
                </div>
                <div class="w3-container" v-else><br><br><br><h2>This Team contains all employees.</h2><br><br><br></div>
                <button class="button" onclick="document.getElementById('addModal').style.display='none'">Cancel</button>

                </div>
              </div>

              <div id="editModal" class="w3-modal">
                <div class="w3-modal-content">
                  <div class="w3-container">

                      <h2>EDIT A TEAM</h2>
                      <br>
                      <label>[NAME] </label>&nbsp;
                      <input type="text" v-model="newTeam.name" placeholder="New name.." required>     
                                
                      <br><br>
                      <label>[DESC] </label>&nbsp;
                      <input type="text" v-model="newTeam.description" placeholder="New Description..." required="required">
                                
                      <br><br>
                      <button class="button" @click="this.editItem(this.team.uid)" onclick="document.getElementById('editModal').style.display='none'">Edit</button>
                      <button class="button" onclick="document.getElementById('editModal').style.display='none'">Cancel</button>   
                  </div>
                </div>
              </div>

              <div id="delModal" class="w3-modal">
              <div class="w3-modal-content">
                <div class="w3-container">
                    <h2>DELETE THIS TEAM</h2>
                    <br>
                    <label>Name:  {{this.team.name}} </label>
                    <br><br>
                    <label>Birth: {{this.team.description}}  </label>
                    <br><br>
                    <button class="button" @click="this.deleteItem(this.team.uid);this.$router.push({name:'Teams'});" onclick="document.getElementById('delModal').style.display='none'">Delete</button>
                    <button class="button" onclick="document.getElementById('delModal').style.display='none'">Cancel</button>      
                </div>
              </div>
              </div> 
        </div>
        <div v-if="list.length!=0">
                <table class="teamList">
                  <thead>
                    <tr>
                      <th class="id">ID</th>
                      <th class="name">JOINED EMPLOYEE</th>
                      <th class="gender">GENDER</th>
                      <th class="dob">BIRTH</th>
                      <th class="option">OPTIONS</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item, index) in list" :key="item.uid">
                      <th class="id">{{++index}}</th>
                      <th class="name">{{item.name}}</th>
                      <th class="gender">{{item.gender}}</th>
                      <th class="dob">{{item.dob}}</th>
                      <th class="option">
                        <button class="button" @click="this.$router.push({name:'Employee', params:{uid: item.uid}})"><i class="fa fa-search"></i> Detail</button>
                        <button class="button" @click="this.chosenEmployee=item;" onclick="document.getElementById('confirmLeaveModal').style.display='block'">
                          <i class="fa fa-remove"></i> Remove</button>
                            <div id="confirmLeaveModal" class="w3-modal">
                              <div class="w3-modal-content">
                                <h2> <b>Remove Employee</b> </h2>
                    
                                <h5>Do you want to remove <b>{{chosenEmployee.name}} from this team?</b></h5>
                                <br>
                                <button class="button" @click="this.leaveTeam(team.uid, chosenEmployee.uid)" 
                                onclick="document.getElementById('confirmLeaveModal').style.display='none';">Confirm</button>
                                <button class="button" onclick="document.getElementById('confirmLeaveModal').style.display='none'">Cancel</button>                            
                              </div>
                            </div>                          
                      </th>
                    </tr>
                  </tbody>

                </table>
               
      </div>
      <div v-else class="teamList"><h1>There is no employee in this team.</h1></div>    
      </div>
    </div>
    </body>
</template>

<script>
import axios from "axios"
export default {
  name: 'Employee',
  data () {
    return {
      uid: null,
      newTeam:{
        name: null,
        description: null,
      },
      team: null,
      employees: [],
      chosenEmployee: {
        uid: null,
        name: null,
        description: null,
      },
      list: [],
      employeeList: [],
    }
  },
  methods: {
    getItem(){
      axios.get(`/api/v1/team/${this.$route.params.uid}`)
      .then((response) => {
        this.team = response.data.data;
      })
    },

    getTeams(){
      axios.get(`/api/v1/employee/`)
      .then((response) => {
        const allTeams = Array.isArray(response.data.data)? response.data.data: [];
        
        axios.get(`/api/v1/team/list/${this.$route.params.uid}`)
        .then((response) => {
          this.employeeList = Array.isArray(response.data.data)? response.data.data: [];
          this.employees = allTeams.filter(item => !this.employeeList.includes(item.uid));
          this.list = allTeams.filter(item => this.employeeList.includes(item.uid));      
        });

      }).catch(error => console.error(error));
    },  

    editItem(uid){
      axios.patch(`/api/v1/team/${uid}`, this.newTeam).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      this.team = this.newTeam;
    },

    deleteItem(uid){
      axios.delete(`/api/v1/team/${uid}`).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      alert("Deleting...")
    },

    addToTeam(tid, eid){
      axios.post(`/api/v1/team/${tid}/employee/${eid}`).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      this.list.push(this.chosenEmployee);
      this.employees = this.employees.filter(item => item.uid!=eid)
    },

    leaveTeam(tid, eid){
      axios.delete(`/api/v1/team/${tid}/employee/${eid}`).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      this.employees.push(this.chosenEmployee);
      this.list = this.list.filter(item => item.uid!=eid);
    },

    update(old){
      this.newTeam.name =new String(old.name);
      this.newTeam.description = new String(old.description);
    }

  },
  mounted() {
    this.getItem();
    this.getTeams();
  },
}
</script>
<style scoped lang="scss">
.info{
    position: relative;
    left: -40%;
}

.all{
   display: flex;
 }
.card {
  box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
  transition: 0.3s;
  width: 20%;
  height: 30%;
  border-radius: 5px;
  margin-top: 5%;
  margin-left: 5%;
}

.card:hover {
  box-shadow: 0 8px 16px 0 rgba(0,0,0,0.2);
}

img {
  border-radius: 5px 5px 0 0;
}

.container {
  padding: 2px 16px;
}

table {
  border: 1px solid green;
  width: 90%;
  margin: auto;
  margin-bottom: 20px;
  margin-top: 50px;
}

table, td, th {
  border: 1px solid green;
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
  width: 20%;
}

.gender{
  width: 10%;
}

.dob{
  width: 15%;
}

.option{
  width: 30%;
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

.add{
  position: absolute;
  top: 15%;
  right: 5%;
}

.teamList{
  position: absolute;
  margin-top: 10%;
  width: 60%;
  margin-left: 10%;
}

input {
  width: 50%;
}
</style>
