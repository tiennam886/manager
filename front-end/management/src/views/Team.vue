<template>
    <body>
        <div class="w3-container">
            <h1>DETAIL</h1>
            <table class="table">
              <tbody v-if="this.list">
              <tr>
                <th>NAME</th>
                <th>{{this.list.name}}</th>
              </tr>
              <tr>
                <th>DESCRIPTION</th>
                <th>{{this.list.description}}</th>
              </tr>
              </tbody>
            </table>
        </div>
    </body>
</template>

<script>
import axios from "axios"
export default {
  name: 'Team',
  data () {
    return {
      list: null,
    }
  },
  methods: {
    getItem(){
      axios.get(`http://localhost:8081/api/v1/team/${this.$route.params.uid}`)
      .then((response) => {
        this.list= response.data.data;
      });
    },
  },
  beforeCreate() {
    console.log(this.list)
    axios.get(`http://localhost:8081/api/v1/team/${this.$route.params.uid}`)
      .then((response) => {
        this.list = response.data.data;
      });
          console.log(this.list)
  }
}
</script>
<style scoped lang="scss">
.info{
    position: relative;
    left: -40%;
}
table, td, th {
  border: 1px solid green;
}

table {
  border: 1px solid green;
  width: 512px;
  height: 216px;
  margin: auto;
}
</style>
