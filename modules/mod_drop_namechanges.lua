function on_pre_precense(event)
    if event and event.origin and event.origin.jitsi_meet_context_user then
        local nick = event.stanza:get_child("nick", "http://jabber.org/protocol/nick")
        local desired_nick = event.origin.jitsi_meet_context_user.name
        
        if nick == nil or nick:get_text() == desired_nick then
           -- ignore, pass on event
           return;
        end
        module:log("info", "====STOPPING CHANGE FROM [%s], REVERTING TO: [%s]", tostring(nick:get_text()), tostring(desired_nick))

        -- Remove the nick tag
        event.stanza:maptags(
           function (tag)
             for k, v in pairs(tag) do
                if k == "name" and v == "nick" then
                   return nil
                end
             end
             return tag
           end
        )

        -- Readd a benign nick tag
        event.stanza:text_tag("nick", desired_nick, { xmlns = "http://jabber.org/protocol/nick" })
    end
end

module:hook("pre-presence/bare", on_pre_precense, 10);
module:hook("pre-presence/full", on_pre_precense, 10);