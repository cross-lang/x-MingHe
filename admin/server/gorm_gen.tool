version: "0.1"
database:
  dsn : "root:123456@tcp(127.0.0.1:3306)/minghe?charset=utf8mb4&parseTime=true&loc=Local"
  db  : "mysql"
  outPath :  "./model/custom"
  onlyModel: true
  modelPkgName  : "custom"
  fieldWithIndexTag : true
  fieldWithTypeTag  : true
  fieldSignable: true
  fieldNullable: true
  tables: [
    "x_user",
    "x_enterprise",
    "x_user_enterprise",
    "x_user_identity_verification"
  ]