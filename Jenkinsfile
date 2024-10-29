pipeline {
  agent any

  tools {
    go '1.23.2'
    nodejs '23.1.0'
  }

  stages {
    stage('Install just') {
      steps {
        sh 'wget https://github.com/casey/just/releases/download/1.36.0/just-1.36.0-x86_64-unknown-linux-musl.tar.gz -O just.tar.gz'
        sh 'tar -zxvf just.tar.gz just'
      }
    }
    stage('Check formatting') {
      steps {
        sh 'cd client && npx prettier --check ./src'
        sh 'test -z $(gofmt -l .)'
      }
    }
    stage('Build') {
      steps {
        sh './just'
      }
    }
    stage('Test go') {
      steps {
        sh 'go test -v ./...'
      }
    }
    stage('SonarQube Analysis') {
      steps {
        script {
          scannerHome = tool 'SonarScanner';
        }
        withSonarQubeEnv('vezgammon') {
          sh "${scannerHome}/bin/sonar-scanner"
        }
      }
    }
  }
}

