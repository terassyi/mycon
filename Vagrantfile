Vagrant.configure("2") do |config|

  config.vm.box = "ubuntu/xenial64"
  config.vm.synced_folder "./", "/usr/local/go/src/github.com/terassyi/mycon"
  config.vm.provision :shell, :path => "./install.sh"

end