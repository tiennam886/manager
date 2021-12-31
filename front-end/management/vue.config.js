module.exports = {
  devServer: {
    proxy: {
      "^/api/v1/employee": {
        target: "http://employee:8082",
        changeOrigin: true,
      },
      "^/api/v1/team": {
        target: "http://team:8081",
        changeOrigin: true,
      },
    },
  },
};
