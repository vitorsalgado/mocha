'use strict'

module.exports = {
  '*.{md,json}': 'prettier --write --ignore-unknown',
  '*.go': ['go fmt', 'go vet'],
}
