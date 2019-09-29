alter table submissions add column ticket_id integer;
--
update submissions set ticket_id=(
  select 
  case 
    when pull_request_id != 0 then pull_request_id
    else null 
  end
  from tasks where submission_id=submissions.id and pull_request_id is not null limit 1);
--
alter table tasks drop column pull_request_id;
