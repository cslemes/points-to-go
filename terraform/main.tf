variable "container_image" {
  type = string
}
variable "project_name" {
  type    = string
  default = "ecs-cluster"
}
variable "aws_region" {
  type    = string
  default = "us-east-2"
}

// points to go
resource "aws_ecs_task_definition" "task_def" {

  family                   = "tmw_services-${var.project_name}"
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
        name  = "USER_DB"
        value = local.db_credentials.username
      },
      {
        name  = "PASSWORD_DB"
        value = local.db_credentials.password
      },
      {
        name      = "DB_HOST"
        valueFrom = data.aws_ssm_parameter.aurora_endpoint.value
      },
      {
        name      = "DB_PORT"
        valueFrom = "3306"
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
    security_groups  = [data.aws_ssm_parameter.sg_tmw_services.value]
    assign_public_ip = true
  }

  service_registries {
    registry_arn   = data.aws_ssm_parameter.service_discovery_service.value
    container_name = "points"
  }
}
