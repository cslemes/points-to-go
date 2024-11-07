variable "container_image" {
  type = string
}
variable "project_name" {
  type    = string
  default = "ecs-cluster"
}


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
  name = format("/%s/ecs/task_role", var.project_name)
}

data "aws_ssm_parameter" "region" {
  name = format("/%s/region", var.project_name)
}

data "aws_ssm_parameter" "mysql_secret" {
  name = format("/%s/mysql_secret", var.project_name)
}


resource "aws_ecs_task_definition" "task_def" {
  family                   = "task_def"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = data.aws_ssm_parameter.ecsTaskExecutionRole.value
  task_role_arn            = data.aws_ssm_parameter.ecs_task_role.value

  container_definitions = jsonencode([{
    name      = "points_to_go"
    image     = var.container_image
    essential = true
    environment = [
      {
        name  = "AWS_REGION"
        value = data.aws_ssm_parameter.region.value
      },
      {
        name  = "SECRET_NAME"
        value = data.aws_ssm_parameter.mysql_secret.name
      }
    ]
    portMappings = [{
      containerPort = 8081
      hostPort      = 8081
      protocol      = "tcp"
    }]
  }])
}

resource "aws_ecs_service" "points_to_go" {
  name            = "points_to_go"
  cluster         = data.aws_ssm_parameter.cluster_id.value
  task_definition = aws_ecs_task_definition.task_def.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [data.aws_ssm_parameter.public_1.value, data.aws_ssm_parameter.public_2.value]
    security_groups  = [""]
    assign_public_ip = true
  }

  service_registries {
    registry_arn   = data.aws_ssm_parameter.service_discovery_service.value
    container_name = "points"
  }
}
