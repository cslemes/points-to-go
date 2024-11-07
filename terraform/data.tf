// Get information from SSM parameters

data "aws_ssm_parameter" "cluster_name" {
  name = format("/%s/ecs/cluster_name", var.project_name)
}

data "aws_ssm_parameter" "cluster_arn" {
  name = format("/%s/ecs/cluster_arn", var.project_name)
}

data "aws_ssm_parameter" "cluster_id" {
  name = format("/%s/ecs/cluster_id", var.project_name)
}

data "aws_ssm_parameter" "public_1" {
  name = format("/%s/vpc/subnet_public_1", var.project_name)
}

data "aws_ssm_parameter" "public_2" {
  name = format("/%s/vpc/subnet_public_2", var.project_name)
}

data "aws_ssm_parameter" "ecsTaskExecutionRole" {
  name = format("/%s/ecs/ecsTaskExecutionRole", var.project_name)
}

data "aws_ssm_parameter" "security_group" {
  name = format("/%s/vpc/security_group", var.project_name)
}

data "aws_ssm_parameter" "service_discovery_namespace" {
  name = format("/%s/ecs/service_discovery_namespace", var.project_name)
}

data "aws_ssm_parameter" "service_discovery_service" {
  name = format("/%s/ecs/service_discovery_service", var.project_name)
}

data "aws_ssm_parameter" "ecs_task_role" {
  name = format("/%s/ecs/ecs_task_role", var.project_name)
}



data "aws_ssm_parameter" "sg_vitess" {
  name = format("/%s/vpc/sg_vitess", var.project_name)

}

data "aws_ssm_parameter" "sg_tmw_services" {
  name = format("/%s/vpc/sg_tmw_services", var.project_name)
}

# data "aws_ssm_parameter" "db_host" {
#   name = format("/%s/vitess/db_host", var.project_name)
# }

# data "aws_ssm_parameter" "db_port" {
#   name = format("/%s/vitess/db_port", var.project_name)
# }
data "aws_secretsmanager_secret" "db_credentials" {
  name = "db_secret"
}

data "aws_secretsmanager_secret_version" "db_credentials_version" {
  secret_id = data.aws_secretsmanager_secret.db_credentials.id
}

locals {
  db_credentials = jsondecode(data.aws_secretsmanager_secret_version.db_credentials_version.secret_string)
}
