FROM ruby:2.2.1

RUN apt-get update && \
      apt-get install -y rsyslog cron && \
      apt-get install -y -q tmux vim ack-grep && \
      apt-get clean

RUN wget https://gist.githubusercontent.com/dimroc/0493d11185ac832ea360/raw/d79c7446c24ef34abd3137f8f3bcbb197cace91d/.bashrc -O /root/.bashrc

# Run Bundle in separate folder to allow better docker caching
# http://ilikestuffblog.com/2014/01/06/how-to-skip-bundle-install-when-deploying-a-rails-app-to-docker/
WORKDIR /tmp
ADD ./Gemfile Gemfile
ADD ./Gemfile.lock Gemfile.lock
RUN gem install bundle
RUN bundle

ADD . /ruby/urbanevents/tools
WORKDIR /ruby/urbanevents/tools

# Add crontab file
ADD crontab /etc/crontab
RUN touch /var/log/cron.log

CMD env > /tmp/.env && rsyslogd && cron && tail -f /var/log/syslog /var/log/cron.log
