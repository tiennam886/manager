<template>
  <div class="table">
    <meta http-equiv="Content-Security-Policy" content="default-src *; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval' http://www.google.com">
    <h1>LIST OF EMPLOYEES</h1>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="https://www.w3schools.com/w3css/4/w3.css">
    
    <button class="buttonAdd" @click="this.renew();" onclick="document.getElementById('addModal').style.display='block'"><i class="fa fa-plus"></i> ADD NEW EMPLOYEE</button>

    <div id="addModal" class="w3-modal">
      <div class="w3-modal-content">
      <div class="w3-container">
          <h2>ADD NEW EMPLOYEE</h2>
          <br>
          <label>Name</label>&nbsp;
          <input type="text" v-model="newEmployee.name" placeholder="Your name.." required="required">
          <br><br>
          <label>Birth</label>&nbsp;
          <input type="date" v-model="newEmployee.dob" required="required">
          <br><br>
          <label>Gender</label>&nbsp;
          <select v-model="newEmployee.gender" required="required"> 
            <option value="male">MALE</option>
            <option value="female">FEMALE</option>
          </select>
          <br><br>
        <button class="button" @click="this.postItem();" onclick="document.getElementById('addModal').style.display='none';">Register</button>
        <button class="button" onclick="document.getElementById('addModal').style.display='none'">Cancel</button>
      </div>
      </div>
    </div>

    <div id="editModal" class="w3-modal">
    <div class="w3-modal-content">
      <div class="w3-container">
          <h2>EDIT AN EMPLOYEE</h2>
          <br>
          <label>[NAME] </label>&nbsp;
          <input type="text" v-model="employee.name" placeholder="New name..." required>

          <br><br>
          <label>[BIRTH]</label>&emsp;&emsp;&emsp;
          <input type="date" v-model="employee.dob" required="required">&emsp;&emsp;&emsp;&emsp;

          <br><br>
          <label>[GENDER]</label>&emsp;&emsp;&emsp;
          <select v-model="employee.gender" required="required"> 
            <option value="male">MALE</option>
            <option value="female">FEMALE</option>
          </select>&emsp;&emsp;&emsp;&emsp;&emsp;

          <br><br>

          <button class="button" @click="this.editItem(this.uid)" onclick="document.getElementById('editModal').style.display='none'">Edit</button>
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
          <th class="gender">GENDER</th>
          <th class="dob">DOB</th>
          <th class="option">OPTIONS</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index)  in list" :key="item.uid">
          <th class="id">{{++index}}</th>
          <th class="name">{{item.name}}</th>
          <th class="gender">{{item.gender}}</th>
          <th class="dob">{{item.dob}}</th>
          <th class="option">
            <button class="button" @click="this.$router.push({name:'Employee', params:{uid: item.uid}})"><i class="fa fa-search"></i> Detail</button>
            <button class="button" @click="this.update(item); this.uid=item.uid;" onclick="document.getElementById('editModal').style.display='block'"><i class="fa fa-edit"></i> Edit</button>
            <button class="button" @click="this.employee=item; this.uid=item.uid" onclick="document.getElementById('delModal').style.display='block'"><i class="fa fa-trash"></i> Delete</button>
          </th>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>


import axios from "axios"
export default {
  name: 'A',  
  data () {
    return {
      list: null,
      employee:{
        name: null,
        gender: null,
        dob: null,
      },
      newEmployee:{
        name: null,
        gender: null,
        dob: null,
      },
      uid: null,
    }
  },
  methods: {
    getItem(){
      axios.get(`${process.env.VUE_APP_EMPLOYEE_URL}/api/v1/employee/`)
      .then((response) => {
        this.list = Array.isArray(response.data.data)? response.data.data: [];
        console.log(this.list);
      }).catch(error => console.error(error));
      

    },

    postItem(){
      axios.post(`${process.env.VUE_APP_EMPLOYEE_URL}/api/v1/employee/`, this.newEmployee).then((res) => {
        this.list.push(res.data.data);
      }).catch(error => console.error(error));
    },

    editItem(uid){
      axios.patch(`${process.env.VUE_APP_EMPLOYEE_URL}/api/v1/employee/${uid}`, this.employee).then((res) => {
        console.log(res)
        location.reload();
      }).catch(error => console.error(error));

    },

    deleteItem(uid){
      axios.delete(`${process.env.VUE_APP_EMPLOYEE_URL}/api/v1/employee/${uid}`).then((res) => {
        this.list = this.list.filter(item => item.uid!=uid)
      console.log(res)}).catch(error => console.error(error));
    },

    renew(){
      this.newEmployee.name =null;
      this.newEmployee.dob =null;
      this.newEmployee.gender =null;
    },

    update(old){
      this.employee.name =new String(old.name);
      this.employee.dob = new String(old.dob);
      this.employee.gender = old.gender
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

thead {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: center;
  background-color: #4CAF50;
  color: white;
}

table {
  border: 1px solid green;
  width: 80%;
  margin: auto;
  margin-bottom: 100px;
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

</style>