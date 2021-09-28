@Library('svs_jenkins-library@develop')_

node('docker-agent'){
  stage('checkOut'){
    cleanWs()
    checkoutRepository()
    }
  stage('build'){
    dockerBuildEnv {
      def appImg = docker.build("${REGISTRY}/${PROJECT_PATH}:latest", "${env.BUILD_OPTS} --build-arg GIT_COMMIT_SHORT=${GIT_COMMIT_SHORT} .")
      appImg.push()
      appImg.push("latest")
    }
  }
}
node('k8s-slave'){
  stage('deploy'){
    checkoutRepository()
    load "kubernetes/env/.kubernetesenv"
    // CHANGE K8S_CREDENTIALS, file
    deployToKubernetes(type: 'deployment', namespace: "${K8S_NAMESPACE}", credentialsId: "${K8S_CREDENTIALS}", file: 'kubernetes/tmpl/deployment.yml.j2', processType: 'jinja2', kubectlVersion: "${KUBECTL_VERSION}")
  }

  stage('clean') {
    cleanWs()
  }
}