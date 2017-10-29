+++
date = "2016-10-30T00:00:00"
description = "Quick Start Guide for Windows."
tags = ["Quick", "Start", "Documentation", "Windows"]
title = "Quick Start for Windows"
author = "jonathan"
+++

This is a development quick start guide when you are using Windows.

1.  Make sure that you have your [GOPATH set up](http://www.wadewegner.com/2014/12/easy-go-programming-setup-for-windows/)

2.  Clone coredns and all dependencies:
      `go get github.com/coredns/coredns`

3.  Navigate to the source:
      `cd $ENV:GOPATH\src\github.com\coredns\coredns`

4.  Fork (but not clone) coredns

5.  Update the origin to point at your repository:
      `git remote set-url origin https://github.com/USERNAME/coredns.git`

6.  Open your editor: `code .`

7.  Create a new file named `Corefile` and populate it:
      ```
      # Only port 53 is supported as NSLOOKUP no longer supports non-standard ports
      .:53 {
          # Your router
          proxy . 192.168.1.1:53
          
          file D:\dev\zone\example.org example.org

          errors
          log
      }
      ```

8.  Create the `example.org` file:
      ```
      example.org.   IN SOA dns.example.org. domains.example.org. (
          2012062701   ; serial
          300          ; refresh
          1800         ; retry
          14400        ; expire
          300 )        ; minimum

      @                        IN NS      dns.example.com.

      @                  42000 IN A       127.0.0.1
      @                  42000 IN A       127.0.0.2
      @                  42000 IN A       127.0.0.3

      api                42000 IN CNAME   sample.service.dns.example.de.
      www                42000 IN CNAME   sample.service.dns.example.de.
      blog               42000 IN CNAME   sample.service.dns.example.de.

      @                   3600 IN MX 1    ASPMX1.L.google.com.
      @                   3600 IN MX 1    ASPMX2.L.google.com.
      @                   3600 IN MX 1    ASPMX3.L.google.com.
      @                    300 IN TXT     "v=spf1 include:_spf.google.com ~all"
      ```

9.  You should be able to execute coredns from VSCode, test it with:
      ```
    > nslookup - localhost
    Default Server:  UnKnown
    Address:  ::1

    > example.org
    Server:  UnKnown
    Address:  ::1

    Name:    example.org
    Addresses:  127.0.0.1
            127.0.0.2
            127.0.0.3
      ```

10.  Use `github.com\coredns\coredns` as though it were your own repository. This is required to ensure that debugging works in VSCode.

11.  You might want to add the following to your global `.gitignore`:

    ```
    coredns.exe
    Corefile
    .vscode
    debug
    ```
