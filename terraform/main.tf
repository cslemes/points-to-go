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
        name  = "AWS_REGION"
        value = var.aws_region
      },
      {
        name  = "SECRET_NAME"
        value = data.aws_ssm_parameter.mysql_secret.value
      },
      {
        name      = "DB_HOST"
        valueFrom = data.aws_ssm_parameter.db_host.value
      },
      {
        name      = "DB_PORT"
        valueFrom = data.aws_ssm_parameter.db_port.value
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

output "points_to_go_public_ip" {
  value = aws_ecs_service.points_to_go.load_balancer[0].container_name
}

// vitess db

resource "aws_ecs_task_definition" "vitess" {
  family                   = "vitess-${var.project_name}"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "512"
  memory                   = "2048"
  execution_role_arn       = data.aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = data.aws_iam_role.ecs_task_role.arn


  container_definitions = jsonencode([
    {
      name      = "vitess"
      image     = "vitess/lite:latest"
      essential = true
      portMappings = [
        {
          containerPort = 15001
          hostPort      = 15001
        },
        {
          containerPort = 3306
          hostPort      = 3306
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "vitess" {
  name            = "vitess-service-${var.project_name}"
  cluster         = data.aws_ssm_parameter.cluster_id.value
  task_definition = aws_ecs_task_definition.vitess.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [data.aws_ssm_parameter.public_1.value, data.aws_ssm_parameter.public_2.value]
    security_groups  = [aws_security_group.vitess.id]
    assign_public_ip = true
  }
}


resource "aws_ssm_parameter" "db_host" {
  name  = "/vitess/db_host"
  type  = "String"
  value = "vitess-service-${var.project_name}"
}

resource "aws_ssm_parameter" "db_port" {
  name  = "/vitess/db_port"
  type  = "String"
  value = "3306"
}
