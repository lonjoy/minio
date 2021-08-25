@Library('jenkinslibrary@master') _

def tools = new org.devops.sonar.tools()

pipeline {
    agent any
    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
    }

    stages {
        stage('Git Checkout') {
            steps {
                script{
                    tools.PrintMsg("获取代码",'green')
                    tools.PrintMsg("构建的模块为: ${params.module}",'green')
                }
                checkout([$class: 'GitSCM',
                branches: [[name: '*/${BRANCH_NAME}']],
                doGenerateSubmoduleConfigurations: false,
                extensions: [], submoduleCfg: [],
                userRemoteConfigs: [[credentialsId: 'b4cfb60e-7489-4964-bc27-8b568d1c4908', url: "${GIT_URL}"]]])
            }
        }

        stage('Start Build wjs-minio') {
            steps {
                script {
                        sh """
                          cd ${WORKSPACE};make build;
                        """
                }
            }
        }

        stage('Docker Build') {
            steps {
                script {
                    sh """
                    cd ${WORKSPACE};
                    tag=date +%Y%m%d%H%M%S
                    docker build  -t registry.cn-shanghai.aliyuncs.com/yineng/wjs-minio:${tag} .
                    """
                }
            }
        }
        stage('Docker push') {
                    steps {
                        script {
                            sh """
                            cd ${WORKSPACE};
                            tag=date +%Y%m%d%H%M%S
                            docker push registry.cn-shanghai.aliyuncs.com/yineng/wjs-minio:${tag}
                            """
                  }
              }
          }
    }
}

