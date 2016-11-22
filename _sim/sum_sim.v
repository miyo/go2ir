module sum_sim();

  reg clk = 1'b0;
  reg reset = 1'b1;
  reg[31:0] counter = 32'h0;
  reg kick = 1'b0;
  wire busy;

  initial begin
    $dumpfile("a.vcd");
    $dumpvars();
  end

  always begin
  #10
  clk <= 1'b0;
  #10
  clk <= 1'b1;
  end

  always @(posedge clk) begin
    counter <= counter + 1;
  end

  always @(posedge clk) begin
    if(counter == 5) reset <= 1'b0;
    if(counter == 30) kick <= 1'b1; else kick <= 1'b0;
    if(counter > 32 && busy == 1'b0) $finish;
  end

   reg signed [32-1 : 0] address;
   reg signed [32-1 : 0] din;
   reg 			 we, oe;

  always @(posedge clk) begin
     if(counter >= 6 && counter < 6+16) begin
	oe <= 1'b1;
	we <= 1'b1;
	address <= (counter - 6);
	din <= (counter - 6);
     end else begin
	we <= 1'b0;
	oe <= 1'b1;
	address <= 32'h0;
	din <= 32'h0;
     end
  end

   wire signed [32-1 : 0] sum_c_din;
   wire signed [32-1 : 0] sum_c_dout;
   wire 		  sum_c_we;
   wire  		  sum_c_oe;
   wire 		  sum_c_empty;
   wire 		  sum_c_full;
   wire signed [32-1 : 0] sum_s_length;
   wire signed [32-1 : 0] sum_s_address_b;
   wire signed [32-1 : 0] sum_s_din_b;
   wire signed [32-1 : 0] sum_s_dout_b;
   wire 		  sum_s_we_b;
   wire 		  sum_s_oe_b;
	 
 sum U(
       .clk(clk), // i
       .reset(reset), // i
       .sum_c_din(sum_c_din), // o
       .sum_c_dout(sum_c_dout), // i
       .sum_c_we(sum_c_we), // o
       .sum_c_oe(sum_c_oe), // o
       .sum_c_empty(sum_c_empty), // i
       .sum_c_full(sum_c_full), // i
       .sum_s_length(sum_s_length), // i
       .sum_s_address_b(sum_s_address_b), // o
       .sum_s_din_b(sum_s_din_b), // o
       .sum_s_dout_b(sum_s_dout_b), // i
       .sum_s_we_b(sum_s_we_b), // o
       .sum_s_oe_b(sum_s_oe_b), // o
       .sum_busy(busy), // o
       .sum_req(kick) // i
       );
   assign sum_c_full = 1'b0;
   assign sum_c_empty = 1'b0;
   

   dualportram#(.DEPTH(4), .WIDTH(32), .WORDS(16)) 
   ram (
	.clk(clk),
	.reset(reset),
   
	.we(we),
	.oe(oe),
	.address(address),
	.din(din),
	.dout(),
     
	.we_b(sum_s_we_b),
	.oe_b(sum_s_oe_b),
	.address_b(sum_s_address_b),
	.din_b(sum_s_din_b),
	.dout_b(sum_s_dout_b),
	
	.length(sum_s_length)
	);

endmodule
