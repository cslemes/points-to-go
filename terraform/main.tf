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

resource "aws_ecs_task_definition" "task_def" {
  family                   = "task_def"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = data.aws_ssm_parameter.ecsTaskExecutionRole.value

  container_definitions = jsonencode([{
    name      = "points_to_go"
    image     = var.container_image
    cpu       = 256
    memory    = 512
    essential = true
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
    registry_arn   = aws_ecs_service.points_to_go.arn
    container_port = 8081
    container_name = "points"
  }
}




