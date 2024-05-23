Pod::Spec.new do |spec|
  spec.name         = 'Gkyle'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/kyleX-network/go-kyleX'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS KyleX Client'
  spec.source       = { :git => 'https://github.com/kyleX-network/go-kyleX.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gkyle.framework'

	spec.prepare_command = <<-CMD
    curl https://gkylestore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gkyle.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
