desc "Run all go files through gofmt"
task :gofmt do
  FileList['**/*.go'].each do |file|
    `gofmt -w #{file} #{file}`
  end
end
