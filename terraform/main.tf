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
  name = format("/%s/network/public_1", var.project_name)
}

data "aws_ssm_parameter" "public_2" {
  name = format("/%s/network/public_2", var.project_name)
}




resource "aws_ecs_task_definition" "task_def" {
  family                   = "task_def"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"

  container_definitions = jsonencode([{
    name      = "my-container"
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

resource "aws_ecs_service" "my_service" {
  name            = "my-service"
  cluster         = data.aws_ssm_parameter.cluster_id
  task_definition = data.aws_ssm_parameter.cluster_arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [data.aws_ssm_parameter.public_1, data.aws_ssm_parameter.public_2] # Replace with your subnet id
    security_groups  = ["default"]                                                        # Replace with your security group id
    assign_public_ip = true
  }
}