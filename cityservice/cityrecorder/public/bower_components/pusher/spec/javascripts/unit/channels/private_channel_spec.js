describe("PrivateChannel", function() {
  var pusher;
  var channel;

  beforeEach(function() {
    pusher = Pusher.Mocks.getPusher({ foo: "bar" });
    channel = new Pusher.PrivateChannel("private-test", pusher);
  });

  describe("after construction", function() {
    it("#subscribed should be false", function() {
      expect(channel.subscribed).toEqual(false);
    });
  });

  describe("#authorize", function() {
    var authorizer;

    beforeEach(function() {
      authorizer = Pusher.Mocks.getAuthorizer();
      spyOn(Pusher.Channel, "Authorizer").andReturn(authorizer);
    });

    it("should create and call an authorizer", function() {
      channel.authorize("1.23", function() {});
      expect(Pusher.Channel.Authorizer.calls.length).toEqual(1);
      expect(Pusher.Channel.Authorizer).toHaveBeenCalledWith(
        channel,
        { foo: "bar" }
      );
    });

    it("should call back with authorization data", function() {
      var callback = jasmine.createSpy("callback");
      channel.authorize("1.23", callback);

      expect(callback).not.toHaveBeenCalled();
      authorizer._callback(false, { foo: "bar" });

      expect(callback).toHaveBeenCalledWith(false, { foo: "bar" });
    });
  });

  describe("#trigger", function() {
    it("should raise an exception if the event name does not start with client-", function() {
      expect(function() {
        channel.trigger("whatever", {});
      }).toThrow(jasmine.any(Pusher.Errors.BadEventName));
    });

    it("should call send_event on connection", function() {
      channel.trigger("client-test", { k: "v" });
      expect(pusher.send_event)
        .toHaveBeenCalledWith("client-test", { k: "v" }, "private-test");
    });

    it("should return true if connection sent the event", function() {
      pusher.send_event.andReturn(true);
      expect(channel.trigger("client-test", {})).toBe(true);
    });

    it("should return false if connection didn't send the event", function() {
      pusher.send_event.andReturn(false);
      expect(channel.trigger("client-test", {})).toBe(false);
    });
  });

  describe("#disconnect", function() {
    it("should set subscribed to false", function() {
      channel.handleEvent("pusher_internal:subscription_succeeded");
      channel.disconnect();
      expect(channel.subscribed).toEqual(false);
    });
  });

  describe("#handleEvent", function() {
    it("should not emit pusher_internal:* events", function() {
      var callback = jasmine.createSpy("callback");
      channel.bind("pusher_internal:test", callback);
      channel.bind_all(callback);

      channel.handleEvent("pusher_internal:test");

      expect(callback).not.toHaveBeenCalled();
    });

    describe("on pusher_internal:subscription_succeded", function() {
      it("should emit pusher:subscription_succeded", function() {
        var callback = jasmine.createSpy("callback");
        channel.bind("pusher:subscription_succeeded", callback);

        channel.handleEvent("pusher_internal:subscription_succeeded", "123");

        expect(callback).toHaveBeenCalledWith("123");
      });

      it("should set #subscribed to true", function() {
        channel.bind(function() {
          expect(channel.subscribed).toEqual(true);
        });
        channel.handleEvent("pusher_internal:subscription_succeeded");
      });
    });

    describe("on other events", function() {
      it("should emit the event", function() {
        var callback = jasmine.createSpy("callback");
        channel.bind("something", callback);

        channel.handleEvent("something", 9);

        expect(callback).toHaveBeenCalledWith(9);
      });
    });
  });
});
