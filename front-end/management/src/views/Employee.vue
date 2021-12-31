<template>
    <body>
      <div class="all" v-if="this.employee">
        <meta http-equiv="Content-Security-Policy" content="default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com">
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
        <div class="card" >
          <div v-if="this.employee.gender ==='male'"><img src="../assets/male_avatar.png" alt="Avatar" style="width:100%"></div>
          <div v-else><img src="../assets/female_avatar2.png" alt="Avatar" style="width:100%"></div>
          <div class="container">
          <h4><b>{{this.employee.name}}</b></h4>  
          <p>{{this.employee.dob}}</p> 
          </div>
          
      </div>

      <div>
        <div class="add">
          <button class="button" onclick="document.getElementById('addModal').style.display='block'"><i class="fa fa-plus"></i> ADD TO NEW TEAM</button>
          <button class="button" @click="this.update(this.employee);" onclick="document.getElementById('editModal').style.display='block'"><i class="fa fa-edit"></i> Edit</button>
          <button class="button" onclick="document.getElementById('delModal').style.display='block'"><i class="fa fa-trash"></i> Delete</button>
              
              <div id="addModal" class="w3-modal" >
                <div class="w3-modal-content">
                <div class="w3-container" v-if="teams.length!=0">
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
                    <tr v-for="(item, index) in teams" :key="item.uid">
                      <th class="id">{{++index}}</th>
                      <th class="name">{{item.name}}</th>
                      <th class="description">{{item.description}}</th>
                      <th class="option">
                        <button class="button" @click="this.chosenTeam=item;" onclick="document.getElementById('confirmAddModal').style.display='block'">
                          <i class="fa fa-plus"></i> Add</button>
                        <div id="confirmAddModal" class="w3-modal">
                          <div class="w3-modal-content">
                            <h2> <b>ADD TO TEAM</b> </h2>
                
                            <h5>Do you want to add this employee to  <b>{{chosenTeam.name}}</b></h5>
                            <br>
                            <button class="button" @click="this.addToTeam(employee.uid, chosenTeam.uid)" 
                            onclick="document.getElementById('confirmAddModal').style.display='none'; document.getElementById('addModal').style.display='none'">Confirm</button>
                            <button class="button" onclick="document.getElementById('confirmAddModal').style.display='none'">Cancel</button>                            
                          </div>
                        </div>
                      </th>
                    </tr>
                  </tbody>
                </table>                 
                </div>
                <div class="w3-container" v-else><br><br><br><h2>This employee is in all teams</h2><br><br><br></div>
                <button class="button" onclick="document.getElementById('addModal').style.display='none'">Cancel</button>

                </div>
              </div>

              <div id="editModal" class="w3-modal">
                <div class="w3-modal-content">
                  <div class="w3-container">

                      <h2>EDIT AN EMPLOYEE</h2>
                      <br>
                      <label>[NAME]</label>&nbsp;
                      <input type="text" v-model="newEmployee.name" placeholder="New name.." required>
                      
                      <br><br>
                      <label>[BIRTH]</label>&nbsp;
                      <input type="date" v-model="newEmployee.dob" required="required">
                                                        
                      <br><br>
                      <label>[GENDER] </label>&nbsp;
                      <select v-model="newEmployee.gender" required="required"> 
                        <option value="male">MALE</option>
                        <option value="female">FEMALE</option>
                      </select>

                      <br><br>

                      <button class="button" @click="this.editItem(this.employee.uid)" onclick="document.getElementById('editModal').style.display='none'">Edit</button>
                      <button class="button" onclick="document.getElementById('editModal').style.display='none'">Cancel</button>
    
                  </div>
                </div>
              </div>

              <div id="delModal" class="w3-modal">
              <div class="w3-modal-content">
                <div class="w3-container">
                    <h2>DELETE AN EMPLOYEE</h2>
                    <br>
                    <label>Name:  {{this.employee.name}} </label>
                    <br><br>
                    <label>Birth: {{this.employee.dob}}  </label>
                    <br><br>
                    <label>Gender: {{this.employee.gender}} </label> 
                    <br><br>
                    <button class="button" @click="this.deleteItem(this.employee.uid);this.$router.push({name:'Employees'});" onclick="document.getElementById('delModal').style.display='none'">Delete</button>
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
                      <th class="name">JOINED TEAM</th>
                      <th class="description">DESCRIPTION</th>
                      <th class="option">OPTIONS</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(item, index) in list" :key="item.uid">
                      <th class="id">{{++index}}</th>
                      <th class="name">{{item.name}}</th>
                      <th class="description">{{item.description}}</th>
                      <th class="option">
                        <button class="button" @click="this.$router.push({name:'Team', params:{uid: item.uid}})"><i class="fa fa-search"></i> Detail</button>
                        <button class="button" @click="this.chosenTeam=item;" onclick="document.getElementById('confirmLeaveModal').style.display='block'">
                          <i class="fa fa-remove"></i> Leave</button>
                            <div id="confirmLeaveModal" class="w3-modal">
                              <div class="w3-modal-content">
                                <h2> <b>LEAVE TEAM</b> </h2>
                    
                                <h5>Do you want to let this employee leave  <b>{{chosenTeam.name}}</b></h5>
                                <br>
                                <button class="button" @click="this.leaveTeam(employee.uid, chosenTeam.uid)" 
                                onclick="document.getElementById('confirmLeaveModal').style.display='none';">Confirm</button>
                                <button class="button" onclick="document.getElementById('confirmLeaveModal').style.display='none'">Cancel</button>                            
                              </div>
                            </div>                          
                      </th>
                    </tr>
                  </tbody>

                </table>
               
      </div>
      <div v-else class="teamList"><h1>This employee is not in any teams.</h1></div>    
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
      newEmployee:{
        name: null,
        gender: null,
        dob: null,
      },
      employee: null,
      teams: [],
      chosenTeam: {
        uid: null,
        name: null,
        description: null,
      },
      list: [],
      teamList: [],
    }
  },
  methods: {
    getItem(){
      axios.get(`/api/v1/employee/${this.$route.params.uid}`)
      .then((response) => {
        this.employee = response.data.data;
      })
    },

    getTeams(){
      axios.get(`/api/v1/team/`)
      .then((response) => {
        const allTeams = Array.isArray(response.data.data)? response.data.data: [];
        
        axios.get(`/api/v1/employee/list/${this.$route.params.uid}`)
        .then((response) => {
          this.teamList = Array.isArray(response.data.data)? response.data.data: [];
          this.teams = allTeams.filter(item => !this.teamList.includes(item.uid));
          this.list = allTeams.filter(item => this.teamList.includes(item.uid));      
        });
  


      }).catch(error => console.error(error));
    },

    editItem(uid){
      axios.patch(`/api/v1/employee/${uid}`, this.newEmployee).then((res) => {
      this.employee = this.newEmployee;  
      console.log(res)
      }).catch(error => console.error(error));
      
    }, 

    deleteItem(uid){
      axios.delete(`/api/v1/employee/${uid}`).then((res) => {
      console.log(res)
      }).catch(error => console.error(error));
      alert("Deleting...");
    },

    addToTeam(eid, tid){
      axios.post(`/api/v1/employee/${eid}/team/${tid}`).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      this.list.push(this.chosenTeam);
      this.teams = this.teams.filter(item => item.uid!=tid)
    },

    leaveTeam(eid, tid){
      axios.delete(`/api/v1/employee/${eid}/team/${tid}`).then((res) => {
      console.log(res)}).catch(error => console.error(error));
      this.teams.push(this.chosenTeam);
      this.list = this.list.filter(item => item.uid!=tid);
    },

    update(old){
      this.newEmployee.name =new String(old.name);
      this.newEmployee.dob = new String(old.dob);
      this.newEmployee.gender = old.gender
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
  width: 10%;
}

.name{
  width: 15%;
}

.description{
  width: 30%;
}

.option{
  width: 25%;
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
</style>
